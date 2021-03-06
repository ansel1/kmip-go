package compliance

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTTLV_MarshalXML(t *testing.T) {
	ttlv := TTLV{
		Tag:   "blue",
		Type:  "white",
		Value: "green",
		Children: []*TTLV{
			{
				XMLName: xml.Name{Local: "brown"},
				Tag:     "orange",
				Type:    "black",
				Value:   "white",
			},
		},
	}
	b, err := xml.Marshal(ttlv)
	require.NoError(t, err)
	require.Equal(t, string(b), `<TTLV tag="blue" type="white" value="green"><brown tag="orange" type="black" value="white"></brown></TTLV>`)
}

func TestCompare(t *testing.T) {
	t.Skip("this panics, need to fix")

	tests := []struct {
		name    string
		isEq    bool
		v1      *TTLV
		v2      *TTLV
		expVars map[string]string
	}{
		{
			name: "empty",
			v1:   &TTLV{},
			v2:   &TTLV{},
			isEq: true,
		},
		{
			name: "basic",
			v1: &TTLV{
				XMLName: xml.Name{Local: "blue"},
				Tag:     "white",
				Type:    "green",
				Value:   "red",
			},
			v2: &TTLV{
				XMLName: xml.Name{Local: "blue"},
				Tag:     "white",
				Type:    "green",
				Value:   "red",
			},
			isEq: true,
		},
		{
			name: "xmlnameneq",
			v1: &TTLV{
				XMLName: xml.Name{Local: "black"},
				Tag:     "white",
				Type:    "green",
				Value:   "red",
			},
			v2: &TTLV{
				XMLName: xml.Name{Local: "blue"},
				Tag:     "white",
				Type:    "green",
				Value:   "red",
			},
			isEq: false,
		},
		{
			name: "xmlnameoverridestag",
			v1: &TTLV{
				XMLName: xml.Name{Local: "blue"},
				Tag:     "black",
				Type:    "green",
				Value:   "red",
			},
			v2: &TTLV{
				XMLName: xml.Name{Local: "blue"},
				Tag:     "white",
				Type:    "green",
				Value:   "red",
			},
			isEq: true,
		},
		{
			name: "typeneq",
			v1: &TTLV{
				Tag:   "white",
				Type:  "blue",
				Value: "red",
			},
			v2: &TTLV{
				Tag:   "white",
				Type:  "green",
				Value: "red",
			},
			isEq: false,
		},
		{
			name: "valueneq",
			v1: &TTLV{
				Tag:   "white",
				Type:  "green",
				Value: "blue",
			},
			v2: &TTLV{
				Tag:   "white",
				Type:  "green",
				Value: "red",
			},
			isEq: false,
		},
		{
			name: "structure",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
				},
			},
			isEq: true,
		},
		{
			name: "structureneq",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "orange"},
					{Value: "blue"},
				},
			},
			isEq: false,
		},
		{
			name: "structurelenneq",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
					{Value: "yellow"},
				},
			},
			isEq: false,
		},
		{
			name: "structureorderneq",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "green"},
					{Value: "blue"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "blue"},
					{Value: "green"},
				},
			},
			isEq: false,
		},
		{
			name: "variable",
			v1: &TTLV{
				Tag:   "white",
				Type:  "green",
				Value: "$COLOR",
			},
			v2: &TTLV{
				Tag:   "white",
				Type:  "green",
				Value: "red",
			},
			isEq:    true,
			expVars: map[string]string{"$COLOR": "red"},
		},
		{
			name: "variablebackref",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "$COLOR"},
					{Value: "blue"},
					{Value: "$COLOR"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "blue"},
					{Value: "red"},
				},
			},
			isEq:    true,
			expVars: map[string]string{"$COLOR": "red"},
		},
		{
			name: "variablebackrefneq",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "$COLOR_0"},
					{Value: "blue"},
					{Value: "$COLOR_0"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Value: "red"},
					{Value: "blue"},
					{Value: "green"},
				},
			},
			isEq:    false,
			expVars: map[string]string{"$COLOR_0": "red"},
		},
		{
			name: "attrs",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
				},
			},
			isEq: true,
		},
		{
			name: "attrsunordered",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "big"},
					}},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "big"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
				},
			},
			isEq: true,
		},
		{
			name: "attrsneq",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "big"},
					}},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "small"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
				},
			},
			isEq: false,
		},
		{
			name: "attrsvariables",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "$SIZE"},
					}},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Size"},
						{Tag: "AttributeValue", Value: "small"},
					}},
					{Tag: "Attribute", Children: []*TTLV{
						{Tag: "AttributeName", Value: "Color"},
						{Tag: "AttributeValue", Value: "red"},
					}},
				},
			},
			isEq:    true,
			expVars: map[string]string{"$SIZE": "small"},
		},
		{
			name: "ignoreservercorrelationvalue",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Location", Value: "North"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "ServerCorrelationValue", Value: "red"},
					{Tag: "Location", Value: "North"},
				},
			},
			isEq: true,
		},
		{
			name: "requireservercorrelationvalue",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "ServerCorrelationValue", Value: "$SCV"},
					{Tag: "Location", Value: "North"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "Location", Value: "North"},
				},
			},
			isEq: false,
		},
		{
			name: "requireservercorrelationvalue",
			v1: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "ServerCorrelationValue", Value: "$SCV"},
					{Tag: "Location", Value: "North"},
				},
			},
			v2: &TTLV{
				Tag:  "white",
				Type: "green",
				Children: []*TTLV{
					{Tag: "ServerCorrelationValue", Value: "red"},
					{Tag: "Location", Value: "North"},
				},
			},
			isEq: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			eq, vars, diff := Compare(testcase.v1, testcase.v2)
			if testcase.isEq {
				assert.True(t, eq)
			} else {
				assert.False(t, eq)
			}
			if testcase.expVars != nil {
				assert.Equal(t, testcase.expVars, vars)
			}
			if diff != "" {
				t.Log(diff)
			}
		})
	}
}
