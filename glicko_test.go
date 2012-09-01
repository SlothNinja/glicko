package glicko

import (
	"testing"
)

type player struct {
        *Params
}

type players []*player

var TestPlayers = []*player {
        &player{NewParams(1500, 200)},
        &player{NewParams(1400, 30)},
        &player{NewParams(1550, 100)},
        &player{NewParams(1700, 300)},
        &player{NewParams(1500, 350)},
}

type varianceTest struct {
        p       *player
        in      Contests
        out     float64
}

func newContest(params *Params, outcome float64) *Contest {
        contest := new(Contest)
        contest.Params = params
        contest.Outcome = outcome
        return contest
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
        &eTest{TestPlayers[0], newContest(TestPlayers[1].Params, 1), 0.6368428359975715},
        &eTest{TestPlayers[0], newContest(TestPlayers[2].Params, 0), 0.4330698170429618},
        &eTest{TestPlayers[0], newContest(TestPlayers[3].Params, 0), 0.3047183664858115},
}

func TestE(t *testing.T) {
	for _, ut := range eTests {
                if e := ut.p.e(ut.in); e != ut.out {
			t.Errorf("TestE() = %+v, want %+v.", e, ut.out)
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
                TestPlayers[0],
                Contests{newContest(TestPlayers[1].Params, 1), newContest(TestPlayers[2].Params, 0), newContest(TestPlayers[3].Params, 0)},
                55477.92847171939,
        },
}

func TestD2(t *testing.T) {
	for _, ut := range d2Tests {
                if d2 := ut.p.d2(ut.in); d2 != ut.out {
			t.Errorf("TestD2() = %+v, want %+v.", d2, ut.out)
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
                TestPlayers[0],
                Contests{newContest(TestPlayers[1].Params, 1), newContest(TestPlayers[2].Params, 0), newContest(TestPlayers[3].Params, 0)},
                1469,
        },
}

func TestRPrime(t *testing.T) {
	for _, ut := range rPrimeTests {
                if rPrime := ut.p.rPrime(ut.in); rPrime != ut.out {
			t.Errorf("TestRPrime() = %+v, want %+v.", rPrime, ut.out)
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
                TestPlayers[0],
                Contests{newContest(TestPlayers[1].Params, 1), newContest(TestPlayers[2].Params, 0), newContest(TestPlayers[3].Params, 0)},
                156,
        },
}

func TestRDPrime(t *testing.T) {
	for _, ut := range rdPrimeTests {
                if rdPrime := ut.p.rdPrime(ut.in); rdPrime != ut.out {
			t.Errorf("TestRDPrime() = %+v, want %+v.", rdPrime, ut.out)
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
                TestPlayers[0],
                Contests{newContest(TestPlayers[1].Params, 1), newContest(TestPlayers[2].Params, 0), newContest(TestPlayers[3].Params, 0)},
                &player{NewParams(1469, 156)},
        },
}

func TestUpdateRating(t *testing.T) {
	for _, ut := range updateRatingTests {
                updatedParams := ut.p.UpdateRating(ut.in)
                if !(updatedParams.rdf == ut.out.Params.rdf && updatedParams.r == ut.out.Params.r) {
			t.Errorf("UpdateRating() = %+v, want %+v.", updatedParams, ut.out.Params)
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
                TestPlayers[1],
                Contests{},
                &player{NewParams(1400, 350)},
        },
}

func TestDecayRating(t *testing.T) {
	for _, ut := range updateDecayTests {
                for i := 0; i < 30; i++ {
                        updatedParams := ut.p.UpdateRating(ut.in)
                        ut.p.Params = updatedParams
                }
                if !(ut.p.rdf == ut.out.Params.rdf && ut.p.r == ut.out.Params.r) {
			t.Errorf("DecayRating() = %+v, want %+v.", ut.p.Params, ut.out.Params)
		}
	}
}
