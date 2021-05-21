module github.com/gravitational/magnet/examples/golang

go 1.13

require (
	github.com/gravitational/magnet v0.2.6
	github.com/gravitational/trace v1.1.11
	github.com/magefile/mage v1.9.0
)

replace (
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)
