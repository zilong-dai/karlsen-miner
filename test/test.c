#include <stdio.h>
#include <stdint.h>

extern void karlsen(uint64_t work[4], uint64_t timestamp, uint64_t nonce, uint64_t res[4], uint8_t log);

int main(int argc, char **argv)
{
    uint64_t work[4] = {0xad0cb5f9b887dcfc, 0xe8e1ff57c9d2c644, 0xf265dff2c6b273f3, 0x41ecb9b85b352b21};
    uint64_t timestamp = 1702373333550;
    uint64_t nonce = 0x85607266d97aea6f;

    uint64_t res[4] = {0, 0, 0, 0};
    uint8_t log = 1;
    karlsen(work, timestamp, nonce, res, log);

    for (int i = 0; i < 4; i++)
    {
        printf("%016lx ", res[i]);
    }
    printf("\n");

    return 0;
}