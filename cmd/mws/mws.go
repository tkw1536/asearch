package main

import (
	"context"
	"fmt"

	"github.com/tkw1536/asearch"
	"github.com/tkw1536/asearch/adapters/mws"
)

const example = `<apply xml:id="p1.1.m1.1.1.cmml" xref="p1.1.m1.1.1">
<csymbol cd="ambiguous" xml:id="p1.1.m1.1.1.1.cmml" xref="p1.1.m1.1.1">superscript</csymbol>
<mws:qvar>x</mws:qvar>
<cn type="integer" xml:id="p1.1.m1.1.1.3.cmml" xref="p1.1.m1.1.1.3">2</cn>
</apply>` // ?x^2

func main() {
	query := mws.MWS{
		URL: "https://ar5search.kwarc.info/api/mws",
	}

	grouper := &asearch.GroupBy[int64, int64]{
		C:    query.Query(context.Background(), example, 1000),
		Kind: func(i int64) int64 { return i },
		Less: func(a, b int64) bool { return a < b },
	}

	results := grouper.Slice(0, 10)
	fmt.Println(results)
}
