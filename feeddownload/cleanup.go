package feeddownload

import (
	"regexp"
	"strings"
)

var (
	invalidFilename    = regexp.MustCompile(`[^-–.0-9A-Za-z-ÁÀȦÂÄǞǍĂĀÃÅǺǼǢĆĊĈČĎḌḐḒÉÈĖÊËĚĔĒẼE̊ẸǴĠĜǦĞG̃ĢĤḤáàȧâäǟǎăāãåǻǽǣćċĉčďḍḑḓéèėêëěĕēẽe̊ẹǵġĝǧğg̃ģĥḥÍÌİÎÏǏĬĪĨỊĴĶǨĹĻĽĿḼM̂M̄ʼNŃN̂ṄN̈ŇN̄ÑŅṊÓÒȮȰÔÖȪǑŎŌÕȬŐỌǾƠíìiîïǐĭīĩịĵķǩĺļľŀḽm̂m̄ŉńn̂ṅn̈ňn̄ñņṋóòôȯȱöȫǒŏōõȭőọǿơP̄ŔŘŖŚŜṠŠȘṢŤȚṬṰÚÙÛÜǓŬŪŨŰŮỤẂẀŴẄÝỲŶŸȲỸŹŻŽẒǮp̄ŕřŗśŝṡšşṣťțṭṱúùûüǔŭūũűůụẃẁŵẅýỳŷÿȳỹźżžẓǯßœŒçÇ]`) // https://stackoverflow.com/questions/22017723/regex-for-umlaut/56293848#56293848
	colonInTitle       = regexp.MustCompile(`\b:`)
	multipleWhitespace = regexp.MustCompile(`\s+`)
)

func FilenameFromTitle(input string) string {
	var output string
	output = strings.ReplaceAll(input, "|", "-")
	output = colonInTitle.ReplaceAllString(output, " - ")
	output = invalidFilename.ReplaceAllString(output, " ")
	output = multipleWhitespace.ReplaceAllString(output, " ")
	output = strings.TrimSpace(output)
	return output
}
