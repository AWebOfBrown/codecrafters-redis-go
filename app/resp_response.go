package main

type RESPResponse interface {
	serialiseRESPTokens() []byte
}

type IndividualRESPResponse struct {
	tokens []*RESPToken
}

type RESPResponseList struct {
	tokens [][]*RESPToken
}

func (irr IndividualRESPResponse) serialiseRESPTokens() []byte {
	var responseData []byte
	for _, tok := range irr.tokens {
		strBytes, ok := tok.Value.([]byte)
		if ok {
			responseData = append(responseData, strBytes...)
		}
	}
	return responseData
}

func NewIndividualRESPResponse(tokens []*RESPToken) *IndividualRESPResponse {
	return &IndividualRESPResponse{
		tokens: tokens,
	}
}

func (rrl RESPResponseList) serialiseRESPTokens() []byte {
	var responseData []byte
	for _, tokenList := range rrl.tokens {
		for _, token := range tokenList {
			strBytes, ok := token.Value.([]byte)
			if ok {
				responseData = append(responseData, strBytes...)
			}
		}
	}
	return responseData
}

func NewRESPResponseList(tokens [][]*RESPToken) *RESPResponseList {
	return &RESPResponseList{
		tokens: tokens,
	}
}
