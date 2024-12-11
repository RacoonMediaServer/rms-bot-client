package command

import "testing"

func TestInterlayerExtras(t *testing.T) {

	i := &Interlayer{}
	InterlayerStore(i, int(5))
	n, ok := InterlayerLoad[int](i)
	if !ok {
		t.Errorf("Unexpected result")
	}
	if n != 5 {
		t.Errorf("Strange value: %d", n)
	}
}
