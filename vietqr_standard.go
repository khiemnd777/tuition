package main

import (
	"errors"
	"fmt"
	"hash/crc32"
	"regexp"
	"strings"

	"github.com/subiz/vietqr"
)

const (
	vietQRPayloadFormat = "01"
	vietQRStatic        = "11"
	vietQRDynamic       = "12"
	napasAID            = "A000000727"
	napasServiceAccount = "QRIBFTTA"
	currencyVND         = "704"
	countryVN           = "VN"
)

type vietQRRequest struct {
	BankBIN       string
	AccountNumber string
	Amount        int
	BillNumber    string
	Purpose       string
	Dynamic       bool
}

func generateVietQR(req vietQRRequest) (string, error) {
	req.BankBIN = onlyDigits(req.BankBIN)
	req.AccountNumber = cleanAccount(req.AccountNumber)
	req.BillNumber = cleanANS(req.BillNumber, 25)
	req.Purpose = cleanANS(req.Purpose, 25)

	if len(req.BankBIN) != 6 {
		return "", errors.New("bank BIN must be 6 digits")
	}
	if _, ok := vietqr.VNBankM[req.BankBIN]; !ok {
		return "", errors.New("bank BIN is not in VietQR bank list")
	}
	if req.AccountNumber == "" {
		return "", errors.New("account number is required")
	}
	if len(req.AccountNumber) > 19 {
		return "", errors.New("account number max length is 19 characters")
	}
	if req.Amount < 0 {
		return "", errors.New("amount must not be negative")
	}

	beneficiary, err := tlv("01",
		mustTLV("00", req.BankBIN)+
			mustTLV("01", req.AccountNumber),
	)
	if err != nil {
		return "", err
	}

	merchantAccount, err := tlv("38",
		mustTLV("00", napasAID)+
			beneficiary+
			mustTLV("02", napasServiceAccount),
	)
	if err != nil {
		return "", err
	}

	method := vietQRStatic
	if req.Dynamic || req.Amount > 0 {
		method = vietQRDynamic
	}

	payload := strings.Builder{}
	payload.WriteString(mustTLV("00", vietQRPayloadFormat))
	payload.WriteString(mustTLV("01", method))
	payload.WriteString(merchantAccount)
	payload.WriteString(mustTLV("53", currencyVND))
	if req.Amount > 0 {
		payload.WriteString(mustTLV("54", fmt.Sprintf("%d", req.Amount)))
	}
	payload.WriteString(mustTLV("58", countryVN))

	additional := buildAdditionalData(req)
	if additional != "" {
		value, err := tlv("62", additional)
		if err != nil {
			return "", err
		}
		payload.WriteString(value)
	}

	withoutCRC := payload.String() + "6304"
	return withoutCRC + crc16CCITTFalse(withoutCRC), nil
}

func buildAdditionalData(req vietQRRequest) string {
	var additional strings.Builder
	if req.BillNumber != "" {
		additional.WriteString(mustTLV("01", req.BillNumber))
	}
	if req.Purpose != "" {
		additional.WriteString(mustTLV("08", req.Purpose))
	}
	return additional.String()
}

func tlv(id, value string) (string, error) {
	if len(id) != 2 || !regexp.MustCompile(`^[0-9]{2}$`).MatchString(id) {
		return "", fmt.Errorf("invalid TLV id %q", id)
	}
	if value == "" {
		return "", fmt.Errorf("empty TLV value for %s", id)
	}
	if len(value) > 99 {
		return "", fmt.Errorf("TLV %s length %d exceeds 99", id, len(value))
	}
	return fmt.Sprintf("%s%02d%s", id, len(value), value), nil
}

func mustTLV(id, value string) string {
	out, err := tlv(id, value)
	if err != nil {
		panic(err)
	}
	return out
}

func crc16CCITTFalse(value string) string {
	var crc uint16 = 0xffff
	for _, b := range []byte(value) {
		crc ^= uint16(b) << 8
		for range 8 {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return fmt.Sprintf("%04X", crc)
}

func cleanANS(value string, maxLen int) string {
	value = vietqr.Ascii(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^A-Za-z0-9 $%*+\-./:]`).ReplaceAllString(value, "")
	value = strings.Join(strings.Fields(value), " ")
	if maxLen > 0 && len(value) > maxLen {
		return value[:maxLen]
	}
	return value
}

func defaultBillNumber(row paymentRow) string {
	parts := []string{
		row.StudentName,
		row.ParentName,
		row.BankBIN,
		row.BankAccount,
		fmt.Sprintf("%d", row.Amount),
		row.Note,
	}
	sum := crc32.ChecksumIEEE([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("SUN%08X", sum)
}
