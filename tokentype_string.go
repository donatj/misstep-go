// Code generated by "stringer -type=TokenType"; DO NOT EDIT

package main

import "fmt"

const _TokenType_name = "TIllegalTEofTOctoHeadingLineTAtSignHeadingLineTExclaimLineTQuestionLineTDashLineTColonLineTAstriskTSignedTWhitespaceTNewLineTStringTEqualsStringTComment"

var _TokenType_index = [...]uint8{0, 8, 12, 28, 46, 58, 71, 80, 90, 98, 105, 116, 124, 131, 144, 152}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
