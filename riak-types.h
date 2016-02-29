struct fetchArgs {
    char* bucketType;
    char* bucket;
    char* key;
};

typedef void (*cb_fn)();
extern void call_tcb_cb(cb_fn);
