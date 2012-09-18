package glicko

import (
	"testing"
)

type player struct {
        r, rd float64
}

func (this *player) R() float64 {
        return this.r
}

func (this *player) RD() float64 {
        return this.rd
}

type players []*player

var testPlayers = []*player {
        &player{1500, 200},
        &player{1400, 30},
        &player{1550, 100},
        &player{1700, 300},
        &player{1500, 350},
}

type varianceTest struct {
        p       *player
        in      Contests
        out     float64
}

func newContest(r, rd, outcome float64) *Contest {
        contest := new(Contest)
        opponent := new(player)
        opponent.r = r
        opponent.rd = rd
        contest.Rated = opponent
        contest.Outcome = outcome
        return contest
}

type rdTest struct {
        in      *player
        out     float64
}

var rdTests = []*rdTest{
        &rdTest{testPlayers[0], 209.74803932337483},
        &rdTest{testPlayers[1], 69.95884504478329},
        &rdTest{testPlayers[2], 118.29725271535261},
}

var defaultOptions = &options{defaultMinRD, defaultMaxRD, defaultT, defaultC}

func TestRD(t *testing.T) {
	for _, ut := range rdTests {
                if newRD := rd(ut.in, defaultOptions); ut.out != newRD {
			t.Errorf("TestRD() = %+v, want %+v.", newRD, ut.out)
		}
	}
}

type gRDTest struct {
        in      float64
        out     float64
}

var gRDTests = []*gRDTest{
        &gRDTest{30, 0.9954980060779407},
        &gRDTest{100, 0.953148974234513},
        &gRDTest{300, 0.7242354637381511},
}

func TestGRD(t *testing.T) {
	for _, ut := range gRDTests {
                if gRD := g(ut.in); ut.out != gRD {
			t.Errorf("TestGRD() = %+v, want %+v.", gRD, ut.out)
		}
	}
}

type eTest struct {
        p       *player
        in      *Contest
        out     float64
}

var eTests = []*eTest{
        &eTest{testPlayers[0], newContest(testPlayers[1].r, testPlayers[1].rd, 1), 0.6369062639623734},
        &eTest{testPlayers[0], newContest(testPlayers[2].r, testPlayers[2].rd, 0), 0.43304012130623676},
        &eTest{testPlayers[0], newContest(testPlayers[3].r, testPlayers[3].rd, 0), 0.304672365112378},
}

func TestE(t *testing.T) {
	for _, ut := range eTests {
                if newE := e(ut.p, ut.in, defaultOptions); newE != ut.out {
			t.Errorf("TestE() = %+v, want %+v.", newE, ut.out)
		}
	}
}

type d2Test struct {
        p       *player
        in      Contests
        out     float64
}

var d2Tests = []*d2Test{
        &d2Test{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].r, testPlayers[1].rd, 1),
                        newContest(testPlayers[2].r, testPlayers[2].rd, 0),
                        newContest(testPlayers[3].r, testPlayers[3].rd, 0),
                },
                55433.47321339554,
        },
}

func TestD2(t *testing.T) {
	for _, ut := range d2Tests {
                if newD2 := d2(ut.p, ut.in, defaultOptions); newD2 != ut.out {
			t.Errorf("TestD2() = %+v, want %+v.", newD2, ut.out)
		}
	}
}

type rPrimeTest struct {
        p       *player
        in      Contests
        out     float64
}

var rPrimeTests = []*rPrimeTest{
        &rPrimeTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].r, testPlayers[1].rd, 1),
                        newContest(testPlayers[2].r, testPlayers[2].rd, 0),
                        newContest(testPlayers[3].r, testPlayers[3].rd, 0),
                },
                1469,
        },
}

func TestRPrime(t *testing.T) {
	for _, ut := range rPrimeTests {
                if newRPrime := rPrime(ut.p, ut.in, defaultOptions); newRPrime != ut.out {
			t.Errorf("TestRPrime() = %+v, want %+v.", newRPrime, ut.out)
		}
	}
}

type rdPrimeTest struct {
        p       *player
        in      Contests
        out     float64
}

var rdPrimeTests = []*rdPrimeTest{
        &rdPrimeTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].r, testPlayers[1].rd, 1),
                        newContest(testPlayers[2].r, testPlayers[2].rd, 0),
                        newContest(testPlayers[3].r, testPlayers[3].rd, 0),
                },
                156,
        },
}

func TestRDPrime(t *testing.T) {
	for _, ut := range rdPrimeTests {
                if newRDPrime := rdPrime(ut.p, ut.in, defaultOptions); newRDPrime != ut.out {
			t.Errorf("TestRDPrime() = %+v, want %+v.", newRDPrime, ut.out)
		}
	}
}

type updateRatingTest struct {
        p       *player
        in      Contests
        out     *player
}

var updateRatingTests = []updateRatingTest{
        updateRatingTest{
                testPlayers[0],
                Contests{
                        newContest(testPlayers[1].r, testPlayers[1].rd, 1),
                        newContest(testPlayers[2].r, testPlayers[2].rd, 0),
                        newContest(testPlayers[3].r, testPlayers[3].rd, 0),
                },
                &player{1469, 156},
        },
}

func TestUpdateRating(t *testing.T) {
	for _, ut := range updateRatingTests {
                newR, newRD, err := UpdateRating(ut.p, ut.in)
                if err != nil {
			t.Errorf("UpdateRating() err = %v, expect nil.", err)
                }
                if !(newRD == ut.out.RD() && newR == ut.out.R()) {
			t.Errorf("UpdateRating() = %v %v, want %v %v.", newR, newRD, ut.out.R(), ut.out.RD())
		}
	}
}

type updateDecayTest struct {
        p       *player
        in      Contests
        out     *player
}

var updateDecayTests = []updateDecayTest{
        updateDecayTest{
                testPlayers[1],
                Contests{},
                &player{1400, 350},
        },
}

func TestDecayRating(t *testing.T) {
	for _, ut := range updateDecayTests {
                for i := 0; i < 31; i++ {
                        var err error
                        ut.p.r, ut.p.rd, err = UpdateRating(ut.p, ut.in)
                        if err != nil {
                                t.Errorf("UpdateRating() err = %v, expect nil.", err)
                        }
                }
                if !(ut.p.rd == ut.out.RD() && ut.p.r == ut.out.R()) {
			t.Errorf("DecayRating() = %v %v, want %v %v.", ut.p.r, ut.p.rd, ut.out.R(), ut.out.RD())
		}
	}
}
