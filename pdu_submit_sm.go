package smpp34

import (
	"bytes"
	"errors"
)

var (
	// Required SubmitSm Fields
	reqSSMFields = []string{
		SERVICE_TYPE,
		SOURCE_ADDR_TON,
		SOURCE_ADDR_NPI,
		SOURCE_ADDR,
		DEST_ADDR_TON,
		DEST_ADDR_NPI,
		DESTINATION_ADDR,
		ESM_CLASS,
		PROTOCOL_ID,
		PRIORITY_FLAG,
		SCHEDULE_DELIVERY_TIME,
		VALIDITY_PERIOD,
		REGISTERED_DELIVERY,
		REPLACE_IF_PRESENT_FLAG,
		DATA_CODING,
		SM_DEFAULT_MSG_ID,
		SM_LENGTH,
		SHORT_MESSAGE,
	}
)

type SubmitSm struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewSubmitSm(hdr *Header, b []byte) (*SubmitSm, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqSSMFields, r)

	if err != nil {
		return nil, err
	}

	s := &SubmitSm{hdr, fields, tlvs}

	return s, nil
}

func (s *SubmitSm) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *SubmitSm) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *SubmitSm) MandatoryFieldsList() []string {
	return reqSSMFields
}

func (s *SubmitSm) GetHeader() *Header {
	return s.Header
}

func (s *SubmitSm) TLVFields() []*TLVField {
	return s.tlvFields
}

func (s *SubmitSm) writeFields() []byte {
	b := []byte{}

	for i, _ := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *SubmitSm) writeTLVFields() []byte {
	b := []byte{}

	for _, v := range s.tlvFields {
		b = append(b, v.Writer()...)
	}

	return b
}

func (s *SubmitSm) Writer() []byte {
	b := append(s.writeFields(), s.writeTLVFields()...)
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(SUBMIT_SM)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)

	return append(h, b...)
}