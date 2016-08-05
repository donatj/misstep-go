package main

import (
	"fmt"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok TokenType // last read token
		lit string    // last read literal
		n   int       // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser
func NewParser(s *Scanner) *Parser {
	return &Parser{
		s: s,
	}
}

func (p *Parser) scan() (tok TokenType, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) unscan() {
	//this is trash
	p.buf.n = 1
}

func (p *Parser) scanIgnoreWhitespace() (tok TokenType, lit string) {
	tok, lit = p.scan()
	if tok == TWhitespace {
		tok, lit = p.scan()
	}
	return
}

func tokenArrayContains(t TokenType, types []TokenType) bool {
	for _, i := range types {
		if i == t {
			return true
		}
	}

	return false
}

// Parse parses scanner result
func (p *Parser) Parse() (*Schema, error) {
	sch := &Schema{}

TableLoop:
	for {
		var tok TokenType
		var lit string
		if tok, lit = p.scanIgnoreWhitespace(); !tokenArrayContains(tok, HeadingTokens) {
			return nil, fmt.Errorf("found %q %s, expected %s", lit, tok, HeadingTokens)
		}

		tbl := &CreateTable{}
		if tok == TAtSignHeadingLine {
			tbl.IsPsuedo = true
		}

		if tok, lit = p.scanIgnoreWhitespace(); tok != TString {
			return nil, fmt.Errorf("found %q %s, expected TableName", lit, tok)
		}

		tbl.TableName = lit
		sch.Tables = append(sch.Tables, tbl)

	TblCommentLoop:
		for {
			comtok, comlit := p.scanIgnoreWhitespace()
			if comtok == TColonLine {
				tbl.TableComment = append(tbl.TableComment, strings.TrimSpace(comlit))
			} else {
				// I don't believe this is nessessary anymore
				// tbl.TableComment = strings.TrimSpace(tbl.TableComment)
				p.unscan()
				break TblCommentLoop
			}
		}

		// ColumnLoop:
		for {
			var coltok TokenType
			var collit string

			if coltok, collit = p.scanIgnoreWhitespace(); !tokenArrayContains(coltok, ColumnTokens) {
				if tokenArrayContains(coltok, HeadingTokens) {
					p.unscan()
					continue TableLoop
				}

				return nil, fmt.Errorf("found %q %s, expected %s", collit, coltok, ColumnTokens)
			}

			col := &TableColumn{}

			switch coltok {
			case TExclaimLine:
				col.ColumnReferenceType = ColumnForeignKeyRegister
			case TQuestionLine:
				col.ColumnReferenceType = ColumnForeignKeyReference
			case TDashLine:
				col.ColumnReferenceType = ColumnRegular
			default:
				return nil, fmt.Errorf("unexpected token: %s", coltok)
			}

			var colntok TokenType
			var colnlit string
			if colntok, colnlit = p.scanIgnoreWhitespace(); colntok != TString {
				return nil, fmt.Errorf("found %q %s, expected Column Name", colnlit, colntok)
			}

			col.ColumnName = colnlit

		ModLoop:
			for {
				modtok, _ := p.scanIgnoreWhitespace()
				switch modtok {
				case TAstrisk:
					col.Nullable = true
				case TSigned:
					col.Signed = true
				case TString:
					p.unscan()
					break ModLoop
				default:
					return nil, fmt.Errorf("unexpected token: %s", modtok)
				}
			} //todo limit to one of each

			_, ctype := p.scan()
			ctype, csize, err := strIntSuffixSplit(ctype)
			if err != nil {
				return nil, err
			}

			cctype, err := ColumnTypes.getColumnType(ctype)
			if err != nil {
				return nil, err
			}

			col.ColumnType = cctype
			col.ColumnSize = csize
			tbl.TableColumns = append(tbl.TableColumns, col)

			for {
				tok, lit := p.scanIgnoreWhitespace()

				fmt.Println(tok, lit)

				if tok != TString && tok != TAstrisk {
					fmt.Println("continuing", tok, lit)
					p.unscan()
					break
				}

				// autoIncr := false
				if tok == TAstrisk {
					// autoIncr = true
					if tbl.AutoIncrColumn == nil {
						tbl.AutoIncrColumn = col
					} else {
						return nil, fmt.Errorf("auto increment column already declared")
					}

					tok, lit := p.scanIgnoreWhitespace()
					if tok != TString {
						return nil, fmt.Errorf("found %q %s, expected %s", lit, tok, TString)
					}
				}

				if lit == "pk" {
				} else if lit[0:1] == "k" {
				} else if lit[0:1] == "u" {
				}
			}

		ColCommentLoop:
			for {
				comtok, comlit := p.scanIgnoreWhitespace()
				if comtok == TColonLine {
					col.ColumnComment = append(col.ColumnComment, strings.TrimSpace(comlit))
				} else {
					// This shouldn't be nessessary
					// col.ColumnComment = strings.TrimSpace(col.ColumnComment)
					p.unscan()
					break ColCommentLoop
				}
			}
			// break
		}

		break
	}

	return sch, nil
}
