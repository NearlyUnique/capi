```bash
capi git create-repo -name=adam

cap <TAB>
git
slack
httpbin

cap git <TAB>
pull
create
pr
issue

cap git pr <TAB>

--name
--from
--todo

cap
    -H 'key: value'      # curl style header
    -k                   # curl style trust insecure sites

    --i-file:filename    # allows variables in other files (.json|.yaml|.xml|.env)
    --i-stdin            # read pipe input, could be automatic

    --b-no-follow        # will default to auto follow 3xx re-directs
    --b-slow-send:rate   # bytes/second
    --b-disconnect:bytes # hard disconnect after x bytes

    --o-format           # (none|json|xml|plain|curl), default is json
    --o-header           # output header values
    --o-no-body          # display body (default display body)
    --o-request          # include HTTP request (includes status)
    --o-status           # status 
    --o-time             # timing details
    --o-file:filename    # with format write to file

key: flags are --i- (input), --b- (behaviour), --o- (output)

```