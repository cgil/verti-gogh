#ifndef RAW_H
#define RAW_H

typedef int(processor_t)(void);

void process_raw(processor_t *p);
void saveraw(char *file);

#endif /* end of include guard: RAW_H */
