package errorx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorSwitch(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		switch testType.NewWithNoMessage().Type() {
		case testType:
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Wrapped", func(t *testing.T) {
		err := testTypeBar1.Wrap(testType.NewWithNoMessage(), "a")
		require.Nil(t, Ignore(err, testTypeBar1))
		require.NotNil(t, Ignore(err, testType))

		switch TypeSwitch(err, testType, testTypeBar1) {
		case testType:
			require.Fail(t, "")
		case testTypeBar1:
			// OK
		case nil:
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Raw", func(t *testing.T) {
		switch TypeSwitch(fmt.Errorf("test non-errorx error"), testType, testTypeBar1) {
		case testType:
			require.Fail(t, "")
		case testTypeBar1:
			require.Fail(t, "")
		case nil:
			require.Fail(t, "")
		default:
			// OK
		}
	})

	t.Run("Nil", func(t *testing.T) {
		switch TypeSwitch(nil, testType, testTypeBar1) {
		case testType:
			require.Fail(t, "")
		case testTypeBar1:
			require.Fail(t, "")
		case nil:
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Supertype", func(t *testing.T) {
		switch TypeSwitch(Decorate(testSubtype0.New("b"), "c"), testType, testTypeBar1) {
		case testTypeBar1:
			require.Fail(t, "")
		case testType:
			// OK
		case nil:
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Subtype", func(t *testing.T) {
		switch TypeSwitch(Decorate(testSubtype0.New("b"), "c"), testSubtype0, testTypeBar1) {
		case testTypeBar1:
			require.Fail(t, "")
		case testType:
			require.Fail(t, "")
		case testSubtype0:
			// OK
		case nil:
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})

	t.Run("SubSubtype", func(t *testing.T) {
		switch TypeSwitch(Decorate(testSubtype0.New("b"), "c"), testSubtype1, testTypeBar1) {
		case testTypeBar1:
			require.Fail(t, "")
		case testType:
			require.Fail(t, "")
		case testSubtype0:
			require.Fail(t, "")
		case testSubtype1:
			require.Fail(t, "")
		case nil:
			require.Fail(t, "")
		default:
			// OK
		}
	})

	t.Run("Ordering", func(t *testing.T) {
		switch TypeSwitch(Decorate(testSubtype0.New("b"), "c"), testSubtype1, testType, testSubtype0, testTypeBar1) {
		case testTypeBar1:
			require.Fail(t, "")
		case testType:
			// OK
		case testSubtype0:
			require.Fail(t, "")
		case testSubtype1:
			require.Fail(t, "")
		case nil:
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})
}

func TestErrorSwitchUnrecognised(t *testing.T) {
	t.Run("Mismatch", func(t *testing.T) {
		switch TypeSwitch(Decorate(testTypeBar2.New("b"), "c"), testTypeBar1) {
		case testTypeBar1:
			require.Fail(t, "")
		case nil:
			require.Fail(t, "")
		case NotRecognisedType():
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Raw", func(t *testing.T) {
		switch TypeSwitch(fmt.Errorf("test"), testTypeBar1) {
		case testType:
			require.Fail(t, "")
		case nil:
			require.Fail(t, "")
		case NotRecognisedType():
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Nil", func(t *testing.T) {
		switch TypeSwitch(nil, testTypeBar1) {
		case testType:
			require.Fail(t, "")
		case nil:
			// OK
		case NotRecognisedType():
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})
}

func TestErrorTraitSwitch(t *testing.T) {
	err := traitTestTimeoutError.Wrap(traitTestError3.NewWithNoMessage(), "a")
	require.True(t, HasTrait(err, Timeout()))
	require.False(t, HasTrait(err, testTrait0))

	t.Run("Wrapped", func(t *testing.T) {
		switch TraitSwitch(err, Timeout(), testTrait0) {
		case testTrait0:
			require.Fail(t, "")
		case Timeout():
			// OK
		case CaseNoError():
			require.Fail(t, "")
		case CaseNoTrait():
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Raw", func(t *testing.T) {
		switch TraitSwitch(fmt.Errorf("test non-errorx error"), Timeout(), testTrait0) {
		case testTrait0:
			require.Fail(t, "")
		case Timeout():
			require.Fail(t, "")
		case CaseNoError():
			require.Fail(t, "")
		case CaseNoTrait():
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Nil", func(t *testing.T) {
		switch TraitSwitch(nil, Timeout(), testTrait0) {
		case testTrait0:
			require.Fail(t, "")
		case Timeout():
			require.Fail(t, "")
		case CaseNoError():
			// OK
		case CaseNoTrait():
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})

	t.Run("NoMatch", func(t *testing.T) {
		switch TraitSwitch(err, testTrait0) {
		case testTrait0:
			require.Fail(t, "")
		case Timeout():
			require.Fail(t, "")
		case CaseNoError():
			require.Fail(t, "")
		case CaseNoTrait():
			// OK
		default:
			require.Fail(t, "")
		}
	})

	t.Run("Ordering", func(t *testing.T) {
		switch TraitSwitch(traitTestTemporaryTimeoutError.Wrap(traitTestError3.NewWithNoMessage(), "a"), Temporary(), Timeout()) {
		case Timeout():
			require.Fail(t, "")
		case Temporary():
			// OK
		case CaseNoError():
			require.Fail(t, "")
		case CaseNoTrait():
			require.Fail(t, "")
		default:
			require.Fail(t, "")
		}
	})
}
