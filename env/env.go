package env

import (
	"github.com/jmorgan1321/SpaceRep/core"
	"text/template"
)

// An Env stores miscellaneous information about our SRS, inorder
// to minimize the need for globals and facilitate testing.
type Env struct {
	SessionLength float32 // Duration of this session in minutes
	Cards         [core.BucketCount][]*core.Card
	Distributions [core.BucketCount]float32
	TmplMap       map[string]*template.Template
}

func New() *Env {
	env := &Env{
		SessionLength: 10.0,
		Cards:         LoadCardsFunc(),
		TmplMap:       map[string]*template.Template{},
	}
	env.Distributions = DistributionFunc(env.Cards, env.SessionLength)

	return env
}

// TODO: remove when we have serialization.
var LoadCardsFunc func() [core.BucketCount][]*core.Card

// DistributionFunc is used by the environment to determine what
// percentage of the time each bucket should display a test to
// the user.
var DistributionFunc func(cards [core.BucketCount][]*core.Card, sessionLength float32) [core.BucketCount]float32

// TODO: comment the rules for how templates get loaded and called.env.go
func (env *Env) LoadTemplates(set string) {
	cardHtmlTmpl1 := `
<div class="front">
    <p>What does the {{.Comp}} <b>{{.Word}}</b> do?</p>
</div>
<div class="back">
    <p><b>{{.Word}}</b> {{.Desc}}</p>
    <img src="{{.Set}}/image/{{.Image}}" height="150" width="150" />
    <p>({{.Hint}})</p>
</div>`
	tmpl1, _ := template.New("test").Parse(cardHtmlTmpl1)
	env.TmplMap[core.WordCard.String()] = tmpl1

	cardHtmlTmpl2 := `
<div class="front">
    <p>What {{.Comp}} {{.Desc}}?</p>
</div>
<div class="back">
    <p><b>{{.Word}}</b></p>
    <img src="{{.Set}}/image/{{.Image}}" height="150" width="150" />
    <p>({{.Hint}})</p>
</div>`
	tmpl2, _ := template.New("test").Parse(cardHtmlTmpl2)
	env.TmplMap[core.DescCard.String()] = tmpl2
}
