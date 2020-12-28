package nationstates_api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNation(t *testing.T) {
	// https://www.nationstates.net/cgi-bin/api.cgi?nation=the_mechalus;q=census+name+flag;scale=46;mode=prank
	xml := `
	<NATION id="the_mechalus">
		<NAME>the Mechalus</NAME>
		<FLAG>https://www.nationstates.net/images/flags/uploads/the_mechalus__47928.png</FLAG>
		<CENSUS>
			<SCALE id="46">
				<PRANK>86</PRANK>
			</SCALE>
		</CENSUS>
	</NATION>`

	nation, err := ParseNation([]byte(xml))
	assert.NoError(t, err)
	assert.Equal(t, "the Mechalus", nation.Name)
	assert.Equal(t, "https://www.nationstates.net/images/flags/uploads/the_mechalus__47928.png", nation.FlagURL)
	assert.Equal(t, 86, nation.GetDefenseForces())
}
