# Multi Target
- Overriding configuration using environment variables
- Creating a hierarchy of build targets (Vertex's within the progressui)
- Error display of a failed target
- Downloading / caching HTTP downloads

`go run mage.go MultipleTargets`

```
❯ go run mage.go MultipleTargets
Logs:    build/logs/latest (build/logs/20200721110651)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 8.5s (5/5) FINISHED
 => CACHED target1                                                                                                            8.5s
 => - target2                                                                                                                 8.5s
 =>   - target3                                                                                                               8.0s
 =>     - ERROR target5                                                                                                       6.0s
 =>   - target4                                                                                                               7.0s
------
 > target5:
#5 0.000 5: hello world:  1
#5 0.500 5: hello world:  2
#5 1.000 Error:  Error on target 5
#5 3.001 Ending
------


❯ ls -l build/logs/latest/
total 20
-rw-r--r-- 1 knisbet knisbet 355 Jul 21 11:07 target1
-rw-r--r-- 1 knisbet knisbet 345 Jul 21 11:07 target2
-rw-r--r-- 1 knisbet knisbet 364 Jul 21 11:07 target3
-rw-r--r-- 1 knisbet knisbet 363 Jul 21 11:07 target4
-rw-r--r-- 1 knisbet knisbet 426 Jul 21 11:07 target5


❯ cat build/logs/latest/target5
Name: target5
Digest: sha256:ad6b8cb41668888a166913dd172726c108e2d4cbd56a0923cdaf944ff690710a
Cached: false
Started: 2020-07-21 11:06:54.285147574 +0000 UTC m=+2.512631085
Completed: <nil>
-----
5: hello world:  1
5: hello world:  2
Error:  Error on target 5
Ending
-----
Vertex: Completed <nil> -> 2020-07-21 11:07:00.285946166 +0000 UTC m=+8.513429620
Vertex: Duration 6.000798535s
Vertex: Error  -> Error on target 5
-----
```

`go run mage.go dl`
```
# First Run
❯ go run mage.go dl
Logs:    build/logs/latest (build/logs/20200721110930)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 5.4s (3/3) FINISHED
 => dl                                                                                                                        5.4s
 => => https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl 44.02MB / 44.02MB            5.4s
 => dep2                                                                                                                      0.7s
 => => http://speedtest-ny.turnkeyinternet.net/100mb.bin 104.86MB / 104.86MB                                                  0.7s
 => dep1                                                                                                                      5.0s
 => => https://speed.hetzner.de/100MB.bin 104.86MB / 104.86MB
 
 # Second Run
 ❯ go run mage.go dl
Logs:    build/logs/latest (build/logs/20200721110956)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 1.0s (3/3) FINISHED
 => dl                                                                                                                        1.0s
 => => https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl                              1.0s
 => dep2                                                                                                                      0.5s
 => => http://speedtest-ny.turnkeyinternet.net/100mb.bin                                                                      0.5s
 => dep1                                                                                                                      0.8s
 => => https://speed.hetzner.de/100MB.bin                                                                                     0.8s


 # Default Cache Directory
 ❯ ls -l build/cache/dl
total 247804
-rw-rw-r-- 1 knisbet knisbet  44023808 Jul 21 11:09 sha256:5c19fb9d14bb061ea1b6156a40d12059017f0b271868e11de36ae5bcc2d031c7
-rw------- 1 knisbet knisbet       117 Jul 21 11:09 sha256:5c19fb9d14bb061ea1b6156a40d12059017f0b271868e11de36ae5bcc2d031c7.wcache
-rw-rw-r-- 1 knisbet knisbet 104857600 Jul 21 11:09 sha256:99c1e95645b4680f618edc8cdf53ce8304db746e58274182a12ea54e91eedaf6
-rw------- 1 knisbet knisbet       101 Jul 21 11:09 sha256:99c1e95645b4680f618edc8cdf53ce8304db746e58274182a12ea54e91eedaf6.wcache
-rw-rw-r-- 1 knisbet knisbet 104857600 Jul 21 11:09 sha256:fc6cb943c2716ade3a5428c96f66a0c3f3be3bcede72a1240aaaf1225e525647
-rw------- 1 knisbet knisbet       101 Jul 21 11:09 sha256:fc6cb943c2716ade3a5428c96f66a0c3f3be3bcede72a1240aaaf1225e525647.wcache

# Etag is used to invalidate cache objects
# sha2sum is recorded and validates the file hasn't changed when returned from the cache
❯ cat build/cache/dl/sha256:5c19fb9d14bb061ea1b6156a40d12059017f0b271868e11de36ae5bcc2d031c7.wcache
etag: '"48e877517c5062f75ead7bf27c776ad1"'
sha2sum: bb16739fcad964c197752200ff89d89aad7b118cb1de5725dc53fe924c40e3f7
 ```


 `go run mage.go dlParallel`
- Starts multiple downloads in parallel, and simulates a failure with one of the downloads

 ```
 ❯ go run mage.go dlParallel
Logs:    build/logs/latest (build/logs/20200721111713)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 6.2s (3/3) FINISHED
 => ERROR downloads                                                                                                           6.2s
 => => http://example.com/non-existant-file                                                                                   6.2s
 => => https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl                              1.0s
 => => http://ipv4.download.thinkbroadband.com/50MB.zip                                                                       1.2s
 => dep2                                                                                                                      0.4s
 => => http://speedtest-ny.turnkeyinternet.net/100mb.bin                                                                      0.4s
 => dep1                                                                                                                      0.8s
 => => https://speed.hetzner.de/100MB.bin                                                                                     0.8s
------
 > downloads:
#1 1.017 url: https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl path: build/cache/dl/sha256:5c19fb9d14bb061ea1b6156a40d12059017f0b271868e11de36ae5bcc2d031c7 error:
#1 1.153 url: http://ipv4.download.thinkbroadband.com/50MB.zip path: build/cache/dl/sha256:f3ba38f992736889d4ec1cc2abfd27de3b5a781ac9de88c74475c82736b99c1f error:
#1 1.154 url: http://example.com/non-existant-file path:  error:
#1 1.154 ERROR REPORT:
#1 1.154 Original Error: *trace.BadParameterError Unexpected status code: 404
#1 1.154 Fields:
#1 1.154   url: http://example.com/non-existant-file
#1 1.154 Stack Trace:
#1 1.154 	/home/knisbet/go/src/github.com/gravitational/magnet/dl.go:88 github.com/gravitational/magnet.(*Magnet).Download
#1 1.154 	/home/knisbet/go/src/github.com/gravitational/magnet/dl.go:35 github.com/gravitational/magnet.(*Magnet).DownloadFuture.func1
#1 1.154 	/usr/lib/go/src/runtime/asm_amd64.s:1373 runtime.goexit
#1 1.154 User Message: Unexpected status code: 404
------
Error: Unexpected status code: 404
exit status 1
 ```