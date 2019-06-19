# Real Go

Highlights of a command-line app development journey, a work in progress.

1. What the app does
   1. read a config
   2. replace values
   3. call http endpoint
   4. output results
2. fork of `pflags` to capture unknown values
3. Concepts and features
   1. Sources - args, env, func-files (TBD), std in (TBD), data files (TBD)
   2. Config readers - json, postman, curl files  (TBD), bash history (TBD)
4. Functions with no side effects, only `main`
   1. references `os.Stdout`
   2. references `os.Stdin`
   3. calls `os.Exit` (prevents `defer`s running)
5. started with functions as the dependency contract, only moved to `interface` when required
6. tab complete
   1. separate package
   2. mock object with tests
   3. debug cli tool to generate inputs for debugging
   4. hook for debugging asa package, packages should not log
7. BDD named tests
    1. avoid table tests except where they _really_ help
    2. hierarchy of tests to prepare scenarios
8. Random stuff
    1. config loader does not allow file read errors to escape. I think a diagnostics hook might be needed
    2. So far I've not got the circle ci test output to display
    3. needs field testing

END
