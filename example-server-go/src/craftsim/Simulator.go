package craftsim

import (
	"fmt"
)

func SimulateSpacecraft() {
	sc_state := LoadCraftJSON("../dictionary.json")

	fmt.Println(sc_state)
}
