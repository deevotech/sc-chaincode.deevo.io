package lib

type Wallet struct {
	KeyType   string
	rsaData   *RSAWallet
	ecdsaData *ECDSAWallet
}

func (w *Wallet) getRSA() *RSAWallet {
	return w.rsaData
}
func (w *Wallet) getECDSA() *ECDSAWallet {
	return w.ecdsaData
}
