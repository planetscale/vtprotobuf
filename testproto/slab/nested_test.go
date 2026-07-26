package slab

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

func makeNestedAddress(i int) *NestedAddress {
	addr := &NestedAddress{
		Street: fmt.Sprintf("%d Main Street, Suite %d", 100+i, i),
		City:   "Springfield",
		Geo:    &NestedGeo{Lat: 47.6 + float64(i)/1000, Lon: -122.3 - float64(i)/1000},
	}
	if i%2 == 0 {
		addr.Region = ptr("WA")
		addr.Zip = ptr(uint32(98100 + i))
		addr.Geo.Accuracy = ptr(float32(i%7) + 0.5)
	}
	return addr
}

func makeNestedLine(i int) *NestedLine {
	line := &NestedLine{
		Sku: fmt.Sprintf("SKU-%08d", i),
		Product: &NestedProduct{
			Title:     fmt.Sprintf("Product %d with a reasonably long title", i),
			Thumbnail: []byte(fmt.Sprintf("thumbnail-bytes-%024d", i)),
			Attributes: map[string]string{
				"color": fmt.Sprintf("color-%d", i%5),
				"size":  fmt.Sprintf("size-%d", i%3),
			},
			Dims: &NestedDimensions{
				Width:    ptr(1.5 * float64(i%10)),
				Height:   ptr(2.5 * float64(i%10)),
				WeightMg: ptr(uint64(i) * 1000),
			},
		},
	}
	if i%2 == 0 {
		line.Quantity = ptr(int64(i%9 + 1))
		line.UnitPrice = ptr(float64(i) * 1.99)
		line.Product.Brand = ptr(fmt.Sprintf("brand-%d", i%4))
		line.Product.Dims.Depth = ptr(0.5 * float64(i%10))
	}
	for d := 0; d < i%3; d++ {
		line.Discounts = append(line.Discounts, &NestedDiscount{
			Code:         fmt.Sprintf("SAVE%d", d),
			AmountMicros: ptr(int64(-1000000 * (d + 1))),
		})
	}
	if i%2 == 0 {
		line.Fulfillment = &NestedLine_Ship{Ship: &NestedShipment{
			Carrier:     "ACME Logistics",
			Destination: makeNestedAddress(i),
			EtaUnix:     ptr(uint64(1700000000 + i)),
		}}
	} else {
		line.Fulfillment = &NestedLine_PickupLocation{PickupLocation: fmt.Sprintf("store-%d", i%10)}
	}
	return line
}

func makeNestedOrder(lines int) *NestedOrder {
	order := &NestedOrder{
		OrderId:  "ORD-2026-000042",
		Sequence: ptr(int64(42)),
		Customer: &NestedParty{
			Name:     "Some Rather Long Customer Name Inc.",
			TaxId:    ptr("TAX-123456789"),
			Billing:  makeNestedAddress(1),
			Delivery: makeNestedAddress(2),
		},
		Totals: &NestedTotals{
			SubtotalMicros: ptr(int64(123450000)),
			TaxMicros:      ptr(int64(12345000)),
			TotalMicros:    ptr(int64(135795000)),
		},
	}
	for i := 0; i < lines; i++ {
		order.Lines = append(order.Lines, makeNestedLine(i))
	}
	if lines > 0 {
		order.Notes = map[string]*NestedNote{}
		for i := 0; i < (lines+3)/4; i++ {
			order.Notes[fmt.Sprintf("note-%d", i)] = &NestedNote{
				Author:    fmt.Sprintf("author-%d", i),
				Text:      "please deliver behind the third flowerpot",
				Timestamp: ptr(int64(1700000000 + i)),
			}
		}
	}
	return order
}

func TestUnmarshalVTSlabNestedEquivalence(t *testing.T) {
	for _, lines := range []int{0, 1, 4, 16, 64} {
		t.Run(fmt.Sprintf("lines=%d", lines), func(t *testing.T) {
			orig := makeNestedOrder(lines)
			data, err := orig.MarshalVT()
			if err != nil {
				t.Fatal(err)
			}
			stock := &NestedOrder{}
			if err := stock.UnmarshalVT(data); err != nil {
				t.Fatal(err)
			}
			slab := &NestedOrder{}
			if err := slab.UnmarshalVTSlab(data); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(orig, slab) {
				t.Fatal("slab-decoded nested message differs from original")
			}
			if !proto.Equal(stock, slab) {
				t.Fatal("slab-decoded nested message differs from stock-decoded one")
			}
		})
	}
}

func TestUnmarshalVTSlabNestedReducesAllocations(t *testing.T) {
	data, err := makeNestedOrder(16).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	stockAllocs := testing.AllocsPerRun(50, func() {
		m := &NestedOrder{}
		if err := m.UnmarshalVT(data); err != nil {
			t.Fatal(err)
		}
	})
	slabAllocs := testing.AllocsPerRun(50, func() {
		m := &NestedOrder{}
		if err := m.UnmarshalVTSlab(data); err != nil {
			t.Fatal(err)
		}
	})
	t.Logf("nested allocs/op: stock=%v slab=%v", stockAllocs, slabAllocs)
	if slabAllocs >= stockAllocs {
		t.Fatalf("slab decoding did not reduce allocations: stock=%v slab=%v", stockAllocs, slabAllocs)
	}
}
