#include <stdlib.h>
#include <stdio.h>
#include <uuid/uuid.h>

int main(void)
{
  uuid_t binuuid;
  uuid_generate_random(binuuid);

  char *uuid = malloc(37);

  #ifndef capitaluuid
  uuid_unparse_upper(binuuid, uuid);

  #elif lowercaseuuid
  uuid_unparse_lower(binuuid, uuid);

  #else
  uuid_unparse(binuuid, uuid);
  #endif

  puts(uuid);
  return 0;
}
