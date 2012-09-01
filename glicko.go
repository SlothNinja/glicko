// Implements the Glicko2 Rating System described at http://www.glicko.net/glicko.html.
package glicko

import . "math"

type Params struct {
        r       float64
        rdf     float64
        minRD   float64
        maxRD   float64
        t       float64
        c       float64
}

// A frequent player may end up with an RD approaching 0 which results in a rating that doesn't change despite improvement
// A minimum RD helps prevent such stagnation.  The value defaults to 30 which is believed to be sufficient for most uses.
func (this *Params) SetMinRD(minRD float64) { this.minRD = minRD }

// RD shows uncertainty in the rating.  A rating for a person that has played games should be more certain than a new player.
// MaxRD provides an upper limit and likely should be set to the RD for a new unranked player.
// Default value is 350.
func (this *Params) SetMaxRD(maxRD float64) { this.maxRD = maxRD }

// T represents the number of rating periods that have past since the present rating was last updated.
// Defaults to 1.
func (this *Params) SetT(t float64) { this.t = t }

// C is a constant used to decay RD over time.
// See, www.glicko.net for a discussion of its value.
// Default is 63.2, which corresponds to a decay length of 30 periods when 350 is used for RD of unranked player.
func (this *Params) SetC(c float64) { this.c = c }

func NewParams(r, rd float64) *Params { return &Params{r, rd, 30, 350, 1, 64} }

// Result records a single game outcome against an opponent and that oppenents Params at the time.
type Contest struct {
        // 1 for win, 0 for loss, 0.5 for tie
        Outcome         float64
        // Params of opponent
        *Params
}

func (this *Params) rd() float64 { return Min(Sqrt(this.rdf * this.rdf + this.t * this.c * this.c), this.maxRD) }

func (this *Params) rd2() float64 { return this.rd() * this.rd() }

// A collection of results for a single player during rating period
type Contests []*Contest

const q = 0.00575646273249 // Log(10) / 400
const q2 = q * q

func (this *Params) rPrime(cs Contests) float64 {
        return Floor(this.r + (q / ((1 / this.rd2()) + (1 / this.d2(cs)))) * this.gsSum(cs))
}

func (this *Params) d2(cs Contests) float64 {
        sum := 0.0
        for _, contest := range cs {
                gRD := g(contest.rd())
                sum += gRD * gRD * this.e1e(contest)
        }
        return 1 / (q2 * sum )
}
func g(rd float64) float64 { return 1 / Sqrt(1 + (3 * q2 * Pow(rd, 2)) / Pow(Pi, 2)) }

func (this *Params) e(c *Contest) float64 { return 1 / (1 + Pow(10, g(c.rd()) * (this.r - c.r) / -400)) }

func (this *Params) e1e(c *Contest) float64 {
        ec := this.e(c)
        return ec * ( 1 - ec)
}

func (this *Params) gsSum(cs Contests) (sum float64) {
        for _, c := range cs {
                sum = g(c.rd()) * (c.Outcome - this.e(c))
        }
        return sum
}

func (this *Params) rdPrime(cs Contests) float64 {
        return Max(Floor(Sqrt(1 / ((1 / this.rd2()) + (1 / this.d2(cs))))), this.minRD)
}

func (this *Params) UpdateRating(cs Contests) *Params {
        if len(cs) == 0 {
                return NewParams(this.r, this.rd())
        }
        return NewParams(this.rPrime(cs), this.rdPrime(cs))
}

func (this *Params) ConfidenceInterval() (float64, float64) {
        twiceRd := 2.0 * this.rdf
        return this.r - twiceRd, this.r + twiceRd
}
