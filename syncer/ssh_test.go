package syncer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testPubKeys1 = `ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEA0KJDLOiiXj9XdMxiCT9KvaKfuxFQi+CIiklaN5hHsNgYOu7TijqyONEu5fONLoAo/cshLa+KuargyTrtizwcP4TPcTXZhhJrM0GUDJragw7SMVIs/5xJBGAyHKJ1YUMGO7+nJTmsCLx6PFOlQYveuriiVVCCZerGCLH+UtSXK3z+l7hx9NiDg3/ylOLc3f3SLxrJKn0gMTgK7BHJFXo4PguuPjWZLVdUDX+XKiqtT2n4IsYs6N9qVFG3zUgNlEjZM47NK/ytAC0max98pK+QNzsuaQOo/IShJ1TOw5wwScflPArVJ2AyROqAe7cfQg7q12I9olASFd3U5NazfZCTYAvWA1kz9UZEWLJ1Br1XOkPqOleMM8KCp/PXzz8H0kISkMIji0/QuiZOPEBsKlszXjlALcXR8Mg1uiZVWy48i9JheyXyj1ToCj6cPScpgFHp3DAGSlKKbE1EFaVfeeyGAnHESlnDDg3Gq5xSsB9Okqm3V5t8GpFaJbV68BxQ4BK6HJ21A3CinV4LdV3hR/OBUbDG2EcI+ZKRDjlpJuu4YU= stace@pretend-machine
ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEAywWhrwq4FjHt+UuwZcZePxtjtZOENFpOjufycaYso2nTlzNwnAQEQRfbqsUxKVtOtGxgApIkUvjRIjNBdJE6iOzvBXZhhJrM0GUDJragw7SMVIs/5xJBGAyHKJ1YUMGO7+nJTmsCLx6PFOlQYveuriiVVCCZerGCLH+UtSXK3z+l7hx9NiDg3/ylOLc3f3SLxrJKn0gMTgK7BHJFXo4PguuPjWZLVdUDX+XKiqtT2n4IsYs6N9qVFG3zUgNlEjZM47NK/ytAC0max98pK+QNzsuaQOo/IShJ1TOw5wwScflPArVJ2AyROqAe7cfQg7q12I9olASFd3U5NazfZCTYAvWA1kz9UZEWLJ1Br1XOkPqOleMM8KCp/PXzz8H0kISkMIji0/QuiZOPEBsKlszXjlALcXR8Mg1uiZVWy48i9JheyXyj1ToCj6cPScpgFHp3DAGSlKKbE1EFaVfeeyGAnHESuXC9wkSeFZCEyMJ+RgJxMkBXNZmyycbwsSqAeGJpMEUDlwzu2GD0obBz0HXqg9J1Xallop5AVDKfeszZcc= stace@another-machine
`
	testPubKeys2 = `
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCqql6MzstZYh1TmWWv11q5O3pISj2ZFl9HgH1JLknLLx44+tXfJ7mIrKNxOOwxIxvcBF8PXSYvobFYEZjGIVCEAjrUzLiIxbyCoxVyle7Q+bqgZ8SeeM8wzytsY+dVGcBxF6N4JS+zVk5eMcV385gG3Y6ON3EG112n6d+SMXY0OEBIcO6x+PnUSGHrSgpBgX7Ks1r7xqFa7heJLLt2wWwkARptX7udSq05paBhcpB0pHtA1Rfz3K2B+ZVIpSDfki9UVKzT8JUmwW6NNzSgxUfQHGwnW7kj4jp4AT0VZk3ADw497M2G/12N0PPB5CnhHf7ovgy6nL1ikrygTKRFmNZISvAcywB9GVqNAVE+ZHDSCuURNsAInVzgYo9xgJDW8wUw2o8U77+xiFxgI5QSZX3Iq7YLMgeksaO4rBJEa54k8m5wEiEE1nUhLuJ0X/vh2xPff6SQ1BL/zkOhvJCACK6Vb15mDOeCSq54Cr7kvS46itMosi/uS66+PujOO+xt/2FWYepz6ZlN70bRly57Q06J+ZJoc9FfBCbCyYH7U/ASsmY095ywPsBo1XQ9PqhnN1/YOorJ068foQDNVpm146mUpILVxmq41Cj55YKHEazXGsdBIbXWhcrRf4G2fJLRcGUr9q8/lERo9oxRm5JFX6TCmj6kmiFqv+Ow9gI0x8GvaQ== demo@test

# comment
`
	testPubKeys = []string{
		"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEA0KJDLOiiXj9XdMxiCT9KvaKfuxFQi+CIiklaN5hHsNgYOu7TijqyONEu5fONLoAo/cshLa+KuargyTrtizwcP4TPcTXZhhJrM0GUDJragw7SMVIs/5xJBGAyHKJ1YUMGO7+nJTmsCLx6PFOlQYveuriiVVCCZerGCLH+UtSXK3z+l7hx9NiDg3/ylOLc3f3SLxrJKn0gMTgK7BHJFXo4PguuPjWZLVdUDX+XKiqtT2n4IsYs6N9qVFG3zUgNlEjZM47NK/ytAC0max98pK+QNzsuaQOo/IShJ1TOw5wwScflPArVJ2AyROqAe7cfQg7q12I9olASFd3U5NazfZCTYAvWA1kz9UZEWLJ1Br1XOkPqOleMM8KCp/PXzz8H0kISkMIji0/QuiZOPEBsKlszXjlALcXR8Mg1uiZVWy48i9JheyXyj1ToCj6cPScpgFHp3DAGSlKKbE1EFaVfeeyGAnHESlnDDg3Gq5xSsB9Okqm3V5t8GpFaJbV68BxQ4BK6HJ21A3CinV4LdV3hR/OBUbDG2EcI+ZKRDjlpJuu4YU= stace@pretend-machine",
		"ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEAywWhrwq4FjHt+UuwZcZePxtjtZOENFpOjufycaYso2nTlzNwnAQEQRfbqsUxKVtOtGxgApIkUvjRIjNBdJE6iOzvBXZhhJrM0GUDJragw7SMVIs/5xJBGAyHKJ1YUMGO7+nJTmsCLx6PFOlQYveuriiVVCCZerGCLH+UtSXK3z+l7hx9NiDg3/ylOLc3f3SLxrJKn0gMTgK7BHJFXo4PguuPjWZLVdUDX+XKiqtT2n4IsYs6N9qVFG3zUgNlEjZM47NK/ytAC0max98pK+QNzsuaQOo/IShJ1TOw5wwScflPArVJ2AyROqAe7cfQg7q12I9olASFd3U5NazfZCTYAvWA1kz9UZEWLJ1Br1XOkPqOleMM8KCp/PXzz8H0kISkMIji0/QuiZOPEBsKlszXjlALcXR8Mg1uiZVWy48i9JheyXyj1ToCj6cPScpgFHp3DAGSlKKbE1EFaVfeeyGAnHESuXC9wkSeFZCEyMJ+RgJxMkBXNZmyycbwsSqAeGJpMEUDlwzu2GD0obBz0HXqg9J1Xallop5AVDKfeszZcc= stace@another-machine",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCqql6MzstZYh1TmWWv11q5O3pISj2ZFl9HgH1JLknLLx44+tXfJ7mIrKNxOOwxIxvcBF8PXSYvobFYEZjGIVCEAjrUzLiIxbyCoxVyle7Q+bqgZ8SeeM8wzytsY+dVGcBxF6N4JS+zVk5eMcV385gG3Y6ON3EG112n6d+SMXY0OEBIcO6x+PnUSGHrSgpBgX7Ks1r7xqFa7heJLLt2wWwkARptX7udSq05paBhcpB0pHtA1Rfz3K2B+ZVIpSDfki9UVKzT8JUmwW6NNzSgxUfQHGwnW7kj4jp4AT0VZk3ADw497M2G/12N0PPB5CnhHf7ovgy6nL1ikrygTKRFmNZISvAcywB9GVqNAVE+ZHDSCuURNsAInVzgYo9xgJDW8wUw2o8U77+xiFxgI5QSZX3Iq7YLMgeksaO4rBJEa54k8m5wEiEE1nUhLuJ0X/vh2xPff6SQ1BL/zkOhvJCACK6Vb15mDOeCSq54Cr7kvS46itMosi/uS66+PujOO+xt/2FWYepz6ZlN70bRly57Q06J+ZJoc9FfBCbCyYH7U/ASsmY095ywPsBo1XQ9PqhnN1/YOorJ068foQDNVpm146mUpILVxmq41Cj55YKHEazXGsdBIbXWhcrRf4G2fJLRcGUr9q8/lERo9oxRm5JFX6TCmj6kmiFqv+Ow9gI0x8GvaQ== demo@test",
	}
)

func TestGetPubKeys(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	err = os.Mkdir(filepath.Join(d, ".ssh"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(filepath.Join(d, ".ssh", "authorized_keys"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(testPubKeys1))
	if err != nil {
		t.Error(err)
	}

	f2, err := os.Create(filepath.Join(d, ".ssh", "authorized_keys2"))
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()
	_, err = f2.Write([]byte(testPubKeys2))
	if err != nil {
		t.Error(err)
	}

	pubkeys, err := getPubKeys(d)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testPubKeys, pubkeys) {
		t.Error(`!reflect.DeepEqual(testPubKeys, pubkeys)`)
	}
}
