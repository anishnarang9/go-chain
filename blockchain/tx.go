package blockchain

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	// if sig matches data theyre authorized to spend
	return in.Sig == data
	// is this output owned by X
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	// verfies that the output was assigned to the person in the first place
	return out.PubKey == data
	// Is new input provide X's authorization
}
