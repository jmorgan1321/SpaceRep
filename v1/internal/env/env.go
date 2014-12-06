package env

import (
	"github.com/jmorgan1321/SpaceRep/v1/internal/builder"
	"github.com/jmorgan1321/SpaceRep/v1/internal/core"
	"text/template"
)

// An Env stores miscellaneous information about our SRS, in order
// to minimize the need for globals and facilitate testing.
type Env struct {
	builder       *builder.Builder
	Cards         core.Deck
	Distributions [core.BucketCount]float32
	TmplMap       map[string]*template.Template
	// SessionLength float32 // Duration of this session in minutes
}

type Option func(*Env)

func New(opts ...Option) *Env {
	env := &Env{}
	// env := &Env{
	//     SessionLength: 10.0,
	//     Cards:         LoadCardsFunc(),
	//     TmplMap:       map[string]*template.Template{},
	// }
	// env.Distributions = DistributionFunc(env.Cards, env.SessionLength)

	for _, opt := range opts {
		opt(env)
	}

	// if env.builder != nil {
	// 	env.Cards = env.builder.LoadDeck()
	// }

	return env
}

func DeckBuilder(d *builder.Builder) Option {
	return func(e *Env) {
		e.builder = d
	}
}

// DistributionFunc is used by the environment to determine what
// percentage of the time each bucket should display a test to
// the user.
// var DistributionFunc f
func DistributionFunc(f func(deck core.Deck, sessionLength float32) [core.BucketCount]float32) Option {
	return func(e *Env) {
		e.Distributions = f(e.Cards, 10.0)
		// e.Distributions = f(env.Cards, env.SessionLength)
	}
}

// func SessionLength(minutes float32) Option {
//     return func(e *Env) {
//         e.SessionLength = minutes
//     }
// }

var tmplIndex = `
<div class="front">
    {{template "front" .}}
</div>
<div class="back">
    {{template "back" .}}
</div>
`

// // TODO: comment the rules for how templates get loaded and called.
// func (env *Env) LoadTemplates(set string) {
//     tmpl, _ := template.New("base").Parse(tmplIndex)
//     filenames := []string{}
//     tmplDir := "html/" + set + "/tmpl"

//     filepath.Walk(tmplDir, func(path string, fi os.FileInfo, err error) error {
//         shouldScan := strings.HasSuffix(path, ".tmpl") && !fi.IsDir()
//         if !shouldScan {
//             return nil
//         }

//         // read contents of template file into string
//         b, _ := ioutil.ReadFile(path)
//         tmplText := string(b)

//         // store tempalte in Template map.
//         filename := fi.Name()[:len(fi.Name())-len(".tmpl")]
//         tmpl2, _ := tmpl.Clone()
//         tmpl2, _ = tmpl2.Parse(tmplText)
//         env.TmplMap[filename] = tmpl
//         filenames = append(filenames, filename)

//         return nil
//     })
// }
