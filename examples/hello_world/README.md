# Hello World
`go run mage.go helloworld`

```
❯ go run mage.go helloworld
Logs:    build/logs/latest (build/logs/20200721105756)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 5.0s (1/1) FINISHED
 => helloworld

 ❯ ls -l build/logs/latest
lrwxrwxrwx 1 knisbet knisbet 14 Jul 21 10:57 build/logs/latest -> 20200721105756

❯ ls -l build/logs/latest/
total 4
-rw-r--r-- 1 knisbet knisbet 483 Jul 21 10:58 helloworld

❯ cat build/logs/latest/helloworld
Name: helloworld
Digest: sha256:936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af
Cached: false
Started: 2020-07-21 10:57:56.635691183 +0000 UTC m=+0.010169618
Completed: <nil>
-----
hello world:  1
hello world:  2
hello world:  3
hello world:  4
hello world:  5
hello world:  6
hello world:  7
hello world:  8
hello world:  9
hello world:  10
-----
Vertex: Completed <nil> -> 2020-07-21 10:58:01.637712355 +0000 UTC m=+5.012190794
Vertex: Duration 5.002021176s
-----
 ```

