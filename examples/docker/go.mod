module github.com/gravitational/magnet/examples/docker

go 1.13

require (
	github.com/gravitational/magnet v0.2.6
	github.com/gravitational/trace v1.1.11
	github.com/magefile/mage v1.9.0
)

replace (
	github.com/gravitational/magnet => ../../
	github.com/jaguilar/vt100 => ../../vendor/github.com/jaguilar/vt100
)
