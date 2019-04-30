# cpdate

cpdate is a "script" that would be used to "copy" all files on one directory to a path tree created from the modified times.

If `cpdate` is used on the following directory:

```
drwxr-xr-x  13 arvos  staff   416B Apr  6 17:40 ..
drwxr-xr-x   3 arvos  staff    96B Apr  6 18:04 .
drwxr-xr-x   3 arvos  staff    96B Apr  6 18:06 cpdate
```

it would create:

```
drwxr-xr-x  13 arvos  staff   416B Apr  6 17:40 ..
drwxr-xr-x   3 arvos  staff    96B Apr  6 18:04 .
drwxr-xr-x   3 arvos  staff    96B Apr  6 18:06 ./2019/05/06/cpdate
drwxr-xr-x   3 arvos  staff    96B Apr  6 18:06 cpdate
```
