package nationstates_api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNation(t *testing.T) {
	// https://www.nationstates.net/cgi-bin/api.cgi?nation=the_mechalus;q=census+fullname+flag+demonym;scale=46;mode=prank
	xml := `
	<NATION id="the_mechalus">
		<FULLNAME>The Empire of the Mechalus</FULLNAME>
		<FLAG>https://www.nationstates.net/images/flags/uploads/the_mechalus__47928.png</FLAG>
		<DEMONYM>Mechalusian</DEMONYM>
		<CENSUS>
			<SCALE id="46">
				<PRANK>86</PRANK>
			</SCALE>
		</CENSUS>
	</NATION>`

	nation, err := ParseNation([]byte(xml))
	assert.NoError(t, err)
	assert.Equal(t, "The Empire of the Mechalus", nation.Name)
	assert.Equal(t, "https://www.nationstates.net/images/flags/uploads/the_mechalus__47928.png", nation.FlagURL)
	assert.Equal(t, "Mechalusian", nation.Demonym)
	assert.Equal(t, 86, nation.GetDefenseForces())
	assert.Equal(t, "https://www.nationstates.net/nation=the_mechalus", nation.GetURL())
}

func TestDefenseCanBeSetAndGet(t *testing.T) {
	nation := Nation{}

	nation.SetDefenseForces(1)
	assert.Equal(t, 1, nation.GetDefenseForces())

	nation.SetDefenseForces(2)
	assert.Equal(t, 2, nation.GetDefenseForces())
}

func TestFlagThumbnailURLsAreGeneratedToMatchTheHostedImageURLs(t *testing.T) {
	mechalus := Nation{
		FlagURL: "https://www.nationstates.net/images/flags/uploads/the_mechalus__47928.png",
	}

	assert.Equal(t, "https://www.nationstates.net/images/flags/uploads/the_mechalus__47928t2.png", mechalus.FlagThumbnailURL())

	eritrea := Nation{
		FlagURL: "https://www.nationstates.net/images/flags/Eritrea.png",
	}

	assert.Equal(t, "https://www.nationstates.net/images/flags/Eritreat2.png", eritrea.FlagThumbnailURL())
}

func TestCanonicalNameLowerCaseAndWithoutSpacesStaysTheSame(t *testing.T) {
	assert.Equal(t, "testlandia", GetCanonicalName("testlandia"))
	assert.Equal(t, "the_mechalus", GetCanonicalName("the_mechalus"))
}

func TestCanonicalNameSpacesReplacedWithUnderscores(t *testing.T) {
	assert.Equal(t, "the_mechalus", GetCanonicalName("the mechalus"))
	assert.Equal(t, "the_west_pacific", GetCanonicalName("The West Pacific"))
}

func TestCanonicalNameUpperCaseChangedToLowerCase(t *testing.T) {
	assert.Equal(t, "the_mechalus", GetCanonicalName("The_Mechalus"))
	assert.Equal(t, "the_west_pacific", GetCanonicalName("The West Pacific"))
}
