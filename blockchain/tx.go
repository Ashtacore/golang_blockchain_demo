package blockchain

type TxInput struct {
	ID        []byte
	Out       int
	Signature string
}

type TxOutput struct {
	Value  int
	PubKey string
}

// Typically there'd be more logic here to match public keys and their signatures
func (in *TxInput) CanUnlockOutputWith(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlockedWith(data string) bool {
	return out.PubKey == data
}
