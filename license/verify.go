package license

import (
	"encoding/json"
	"errors"
	"github.com/hyperboloide/lk"
	"github.com/toolkits/file"
	"os"
	"time"
)

//verifyLicense, to verify license.
func VerifyLicense() (end bool, err error) {
	licenseB32, e := file.ToString(".license.dat")
	if err != nil {
		err = e
		return
	}
	// a base32 encoded private key generated by lkgen gen
	const privateKeyBase32 = "FD7YCAYBAEFXA22DN5XHIYLJNZSXEAP7QIAACAQBANIHKYQBBIAACAKEAH7YIAAAAAFP7AYFAEBP7BQAAAAP7GP7QIAWCBGNHH2WMESVKT5R564LJ3SVQGSSO3W6S6RXOMLI5PZAQJU463QTQWSJZ4M76QNZYUOI6SD3HLHWDJZ7N66DEXD6JMZWEF52QV2ZHE5JNH36M2MP2V24KY7SNUBQZQPPCBWGK2AEHQONXQKKSYCLPW4YWGLVAEYQERU2D3ML3WGSBQRS5RZ3EWGNBWF3DXU4U4FNGV4X7DRT265MNRIW3LVTHP3IY5DYKYPEZF5G7SFGYQAA===="

	// Unmarshal the private key
	privateKey, e := lk.PrivateKeyFromB32String(privateKeyBase32)
	if err != nil {
		err = e
		return
	}

	publicKey := privateKey.GetPublicKey()

	// Unmarshal the customer license.
	license, e := lk.LicenseFromB32String(licenseB32)
	if err != nil {
		err = e
		return
	}

	// validate the license signature.
	if ok, e := license.Verify(publicKey); err != nil {
		err = e
		return
	} else if !ok {
		err = errors.New("Invalid license signature")
		return
	}

	result := struct {
		End time.Time `json:"end"`
	}{}

	// unmarshal the document.
	if e := json.Unmarshal(license.Data, &result); e != nil {
		err = e
		return
	}

	// Now you just have to check that the end date is after time.Now() then you can continue!
	if result.End.After(time.Now()) {
		return true, nil
	}
	return false, nil
}

// CheckLicenseExpire 1 hour.
func CheckLicenseExpire() {
	for {
		if end, err := VerifyLicense(); err != nil || end == false {
			os.Exit(0)
		}
		time.Sleep(time.Hour)
	}
}
