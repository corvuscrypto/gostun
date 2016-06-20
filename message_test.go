package gostun

import "testing"

var testTID = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

func TestMessageMarshaling(T *testing.T) {
	//create the mock Message without Attributes
	msg := new(Message)
	msg.Magic = StunMagic
	msg.MessageLength = 0
	msg.MessageType = BindingRequest
	msg.TID = testTID
	expected := []byte{0, 1, 0, 0, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	actual, _ := Marshal(msg)
	if string(actual) != string(expected) {
		T.Errorf("Expected to receive %v, but instead received %v", expected, actual)
	}

	//add some Attributes
	msg.Attributes = make(map[uint16][]byte)
	//should have no padding
	msg.Attributes[AttributeChangeAddress] = []byte{0, 1, 2, 3}
	expected = []byte{0, 1, 0, 8, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 3, 0, 4,
		0, 1, 2, 3}
	actual, _ = Marshal(msg)
	if string(actual) != string(expected) {
		T.Errorf("Expected to receive %v, but instead received %v", expected, actual)
	}
	//should have padding of 2
	msg.Attributes[AttributeChangeAddress] = []byte{0, 1, 2, 3, 2, 2}
	expected = []byte{0, 1, 0, 12, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 3, 0, 6,
		0, 1, 2, 3, 2, 2, 0, 0}
	actual, _ = Marshal(msg)
	if string(actual) != string(expected) {
		T.Errorf("Expected to receive %v, but instead received %v", expected, actual)
	}
}

func TestMessageUnMarshaling(T *testing.T) {
	data := []byte{0, 1, 0, 0, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	msg, _ := UnMarshal(data)
	if msg.Attributes != nil {
		T.Errorf("Encountered Attributes where there shouldn't have been")
	}
	if msg.Magic != StunMagic {
		T.Errorf("Incorrect Magic Value")
	}
	if msg.MessageLength != 0 {
		T.Errorf("Incorrect Message Length")
	}
	if msg.MessageType != BindingRequest {
		T.Errorf("Incorrect Message Type")
	}
	if string(msg.TID) != string(testTID) {
		T.Errorf("Incorrect TID")
	}

	//with AttributeChangeAddress added (no padding)
	data = []byte{0, 1, 0, 8, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 3, 0, 4,
		0, 1, 2, 3}
	msg, _ = UnMarshal(data)
	if string(msg.Attributes[AttributeChangeAddress]) != string([]byte{0, 1, 2, 3}) {
		T.Errorf("Incorrect Attribute value")
	}
	if msg.Magic != StunMagic {
		T.Errorf("Incorrect Magic Value")
	}
	if msg.MessageLength != 8 {
		T.Errorf("Incorrect Message Length")
	}
	if msg.MessageType != BindingRequest {
		T.Errorf("Incorrect Message Type")
	}
	if string(msg.TID) != string(testTID) {
		T.Errorf("Incorrect TID")
	}

	//with AttributeChangeAddress added (2 byte padding)
	data = []byte{0, 1, 0, 12, 33, 18, 164, 66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 3, 0, 6,
		0, 1, 2, 3, 2, 2, 0, 0}
	msg, _ = UnMarshal(data)
	if string(msg.Attributes[AttributeChangeAddress]) != string([]byte{0, 1, 2, 3, 2, 2}) {
		T.Errorf("Incorrect Attribute value")
	}
	if msg.Magic != StunMagic {
		T.Errorf("Incorrect Magic Value")
	}
	if msg.MessageLength != 12 {
		T.Errorf("Incorrect Message Length")
	}
	if msg.MessageType != BindingRequest {
		T.Errorf("Incorrect Message Type")
	}
	if string(msg.TID) != string(testTID) {
		T.Errorf("Incorrect TID")
	}

}
