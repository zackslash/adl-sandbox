package main

import (
	"fmt"

	"github.com/cubex/potens-go/adl"
)

// runSandbox is executed every 5 seconds (ADL sandbox code goes in this method)
func runSandbox(count int) {
	ent1 := app.ADL("THIS-IS-A-FID")
	/*err := ent1.Read(adl.Counter("propX"))
	if err != nil {
		fmt.Printf(err.Error())
	}*/

	// Test counters
	testStrStart := "this is data!"
	ent1.Write("propX", testStrStart)

	// Test Meta
	ent1.WriteMeta("propX", "this is meta")

	// Test sets
	ent1.AddSetItem("propX", "test1")
	ent1.AddSetItem("propX", "test2")
	ent1.AddSetItem("propX", "test3")
	//ent1.RemoveSetItem("propX", "test2")

	// Test Lists
	testListName := "TESTLIST"
	ent1.AddListItem(testListName, "1", "ONE")
	ent1.AddListItem(testListName, "2", "TWO")
	ent1.AddListItem(testListName, "3", "THREE")
	ent1.AddListItem(testListName, "4", "FOUR")

	//ent1.IncrementCounter("propX")
	ent1.Write("propY", "This is Y data")
	ent1.Commit()

	//adl.Counter("propX"),
	//adl.PropertiesWithPrefix("prop"),
	//adl.ListRange(testListName, "1", "", 0)
	//adl.Property("propX")

	ent1.Read(adl.PropertiesWithPrefix("prop"), adl.Meta("propX"), adl.Set("propX"))
	dataA := ent1.Get("propX")
	//countA = ent1.GetCounter("propX")
	data3 := ent1.GetSet("propX")
	data4 := ent1.Get("propY")
	meta1 := ent1.GetMeta("propX")
	//list := ent1.GetList(testListName)

	fmt.Printf("\nItem:%s-%s-%s\n", dataA, data4, meta1)
	//fmt.Printf("Counter:%d\n", countA)

	for _, d := range data3 {
		fmt.Printf("SET-ITEM:%s\n", d)
	}

	if len(data3) == 0 {
		fmt.Printf("No set items returned\n")
	}

	/*for _, d := range list {
		fmt.Printf("LIST-ITEM:%s - %s\n", d.Key, d.Value)
	}

	if len(list) == 0 {
		fmt.Printf("No list items returned\n")
	}*/
}
