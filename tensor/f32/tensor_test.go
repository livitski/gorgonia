package tensorf32

import (
	"testing"

	types "github.com/chewxy/gorgonia/tensor/types"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	t.Log("Standard, expected way of creating an ndarray")
	backingGood := make([]float32, 2*2*6)
	T := NewTensor(WithShape(2, 2, 6), WithBacking(backingGood))

	expectedStrides := []int{12, 6, 1}
	assert.Equal(expectedStrides, T.Strides(), "Unequal strides")

	expectedDims := 3
	assert.Equal(expectedDims, T.Dims(), "Unequal dims")

	t.Log("Creating with just passing in a backing")
	T = NewTensor(WithBacking(backingGood)) // if you do this in real life without specifying a shap, you're an idiot
	expectedShape := types.Shape{len(backingGood)}
	assert.Equal(expectedShape, T.Shape(), "Unequal shape")

	t.Log("Creating with just a shape")
	T = NewTensor(WithShape(1, 3, 5))
	assert.Equal(15, T.Size(), "Unequal size")

	t.Log("Creating an ndarray with a mis match shape and elements")
	backingBad := []float32{1, 2, 3, 4}
	badBackingF := func() {
		NewTensor(WithBacking(backingBad), WithShape(2, 2, 6))
	}
	assert.Panics(badBackingF, "Calling NewNDArray with bad backing should have panick'd")

	t.Logf("Making a scalar value a Tensor")
	T = NewTensor(AsScalar(3.1415))
	assert.Equal(0, len(T.Shape()), "Expected a 1D shape")

	t.Log("Creating a ndarray with nothing passed in")
	noshapeF := func() {
		NewTensor()
	}
	assert.Panics(noshapeF, "Calling NewNDArray() without a shape should have panick'd")

}

func TestReshape(t *testing.T) {
	assert := assert.New(t)
	var T *Tensor
	var backing []float32
	var err error

	t.Log("Testing standard reshape")
	backing = make([]float32, 2*2*6)
	T = NewTensor(WithShape(2, 2, 6), WithBacking(backing))
	if err = T.Reshape(12, 2); err != nil {
		t.Errorf("There should be no error. Got %v instead", err)
	}

	expectedShape := types.Shape{12, 2}
	assert.Equal(expectedShape, T.Shape(), "Unequal shape")

	t.Log("Testing wrong reshape")
	if err = T.Reshape(12, 3); err == nil {
		t.Errorf("There should have been an error")
	}
}

func TestOnes(t *testing.T) {
	assert := assert.New(t)
	var T *Tensor
	var backing []float32
	// var err error

	t.Log("Testing usual use case")
	backing = []float32{1, 1, 1, 1}
	T = Ones(2, 2)

	expectedShape := types.Shape{2, 2}
	assert.Equal(expectedShape, T.Shape())
	assert.Equal(backing, T.data)

	t.Log("Testing stupid sizes: no size")
	T = Ones()
	assert.Nil(T.Shape())
	assert.Equal([]float32{1}, T.data)
}

func TestClone(t *testing.T) {
	assert := assert.New(t)

	backing := []float32{1, 2, 3, 4, 5, 6}
	T := NewTensor(WithBacking(backing), WithShape(2, 3))

	T1000 := T.Clone()
	// make sure that they are two different pointers, or else funny corruptions might happen
	if T.AP == T1000.AP {
		t.Error("Access Patterns must be two different objects")
	}
	// BUT the value must be the same
	assert.EqualValues(T.AP, T1000.AP)
	assert.Equal(T.data, T1000.data)
}
