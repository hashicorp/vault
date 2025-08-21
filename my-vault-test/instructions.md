start vault:

./bin/vault server -dev -dev-root-token-id=root

some usage of common endpoints:

1) Enable Transit (API)
curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
  -X POST -d '{"type":"transit"}' \
  "$VAULT_ADDR/v1/sys/mounts/transit"

2) Create a key (simulates your future “Kyber” key, but using AES now)
curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
  -X POST -d '{"type":"kyber768"}' \
  "$VAULT_ADDR/v1/transit/keys/my-kyber22"

3) Encrypt (same body shape you’ll keep for Kyber)
PLAINTEXT_B64=$(printf 'Hello Kyber768' | base64 | tr -d '\n')

CT=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
  -X POST -d "{\"plaintext\":\"$PLAINTEXT_B64\"}" \
  "$VAULT_ADDR/v1/transit/encrypt/my-kyber22" \
  | jq -r .data.ciphertext)

echo "$CT"

4) Decrypt
PT_B64=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
  -X POST -d "{\"ciphertext\":\"$CT\"}" \
  "$VAULT_ADDR/v1/transit/decrypt/my-kyber22" \
  | jq -r .data.plaintext)

# macOS decode:
printf "%s" "$PT_B64" | base64 -D; printf "\n"

testing:
make test // runs tests for all components, but the docker image isn't native to linux and fails.

make test TEST=./builtin/...
make test TEST=./builtin/logical/transit/...


Testing the integration:

create key and run flow:
izi0511@Igors-MacBook-Pro vault % curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
-X POST -d '{"type":"kyber768"}' \
"$VAULT_ADDR/v1/transit/keys/my-kyber22"
{"request_id":"b722afe5-ee3c-4942-09f6-736656a24c57","lease_id":"","renewable":false,"lease_duration":0,"data":{"allow_plaintext_backup":false,"auto_rotate_period":0,"deletion_allowed":false,"derived":false,"exportable":false,"imported_key":false,"keys":{"1":{"certificate_chain":"","creation_time":"2025-08-21T09:18:32.732955+03:00","hybrid_public_key":"","name":"crystals-kyber","public_key":"omFqJllRBOIHKDsq+AQFlxKgiQum73cnivSIKtKE5RofORZrk3deqWWnuNJ2UsSPrzdCa+LDDANsaQcM1piRqmGHpvSqZFEgw2jCo2LK0xuSJ+XMH3O14jhpJCK9nvSEjJuJztPH1bQgveIqJiFji+qyBdZRo5UUtShOgWs/trkv/KtknGtGdwPIPKcySbwlxeRwvvySWpBjZgAQpBa9+VicGqeY9He2zVMndqZszWMWZksgHnej5MuyUfIpqvchr5AMueA4ZZJFFqc+ErVXT2TFqVFOAqqSoQDQB/OF7SmbGvhl80oNA3GOHigKhwVBKPxjRxaIdPt9SckrI7K1Fis13ty0aACLFzIsRJgBvAlK/lRh7Atra6QHm3en3ZJCQKbAtBmbboSfTgkKDhqacxwYs3fKCXJ0kVpxU/yQr+Mn5yF3tCes5HfNa1VGEunDXGSyTtYyIMylWphhTOAlHVE0Z5KD5BQ0kxlc9DGfQvnDBEuOazKFQhKqyhgq0LsEZjwLHAS7giYacJEGHxCbUVd9HFbHAvBy5yl1q/wv1ZpVsnE6DCSqkeIRwsicT/dUyeQjF+LJIum7kBMSrQRakiZQ9zyI4vKW+bSc76CelLE1wYWYySena0YfLPAwJDVvygY66JqJ8dK8fPdFKYuz+oN1cKANk4CYFUBmfyEG50syMRpwNLJ5JsNGAWu+mtggprkyssBTM0Y+0AJdw4eu01gCnoVNABy4rSN0+JIpm1UhVyofKpprShmQNaGbk4VYw/wIFUAyoaM9gnVDb7lANahzO5lz5MmKAwCJWNlvqlzC1pYRfUiPUlmh8Jx6kjSKvDISvhJk45grdNNjwgcc7eS4ZVtHSxC50LJa/UdKGjqXHsEB/MATRIQQ1YAaJXeW9/MZzEgbDqaMiPks+jYMO0JrbGcv2QFX3lWSjCaj7lHDQmLPPIGsNUNB28ukJrzNOPCO7bSvr7x80ELNLKevx0hXekRxI1F8jqeMUbt+SGzDF3q5P/ajm6EGuDzIwWUqjRkMzYZML7eiZilZFgXHYqxLBZOUwkUq59qRMzlnsZGqWIsxEikRwNpmDoFuI9ZYswFD2gA2GeUjm9eKsEFxmmWpdWymX1UMwMwJpTQN8iHBLmeh+oZc4ie9m0uqS7mfQMEMSdDGDdG56lRMhXuHMSahnPtafTImPQK1/POOI9oizcScx/SF+6muKvslZJSh34xvZqU+jyqnOcAlYHYCwPO4JEhVAnujXJZvZoyY0fZOgiOMuuFTkQFnaNorYQEzixOm4qpP4qunr+p9Zgon78K1BOk/w7aMpBpxLfa/aKlh/9hP5YR3qXUfh4mEQtuYBABZcgmFU5ZgAmOk3NNvOzwEpMLKQzpiYhpHKzAK9ssjy/fKGjkgj5yX73jEz4WjHkMJYjaFgXOvb/Im2dWPK5kDndjIrBIrY+eHnbeZvLU4tWuszvV6sqg+hKgksyJIbSAzq8xuW3a0tfSbdAK6+VUOdhAajwlOyvcDgCk+xslx6FoLEiwNTwG4Pda98OwnHOGsoLlVXDDGz/tcIUxdBhXfNg5rXxPLFCJ4cEMYmdU="}},"latest_version":1,"min_available_version":0,"min_decryption_version":1,"min_encryption_version":0,"name":"my-kyber22","supports_decryption":true,"supports_derivation":false,"supports_encryption":true,"supports_signing":false,"type":"kyber768"},"wrap_info":null,"warnings":null,"auth":null,"mount_type":"transit"}

izi0511@Igors-MacBook-Pro vault % PLAINTEXT_B64=$(printf 'Hello Kyber768' | base64 | tr -d '\n')

izi0511@Igors-MacBook-Pro vault % CT=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
-X POST -d "{\"plaintext\":\"$PLAINTEXT_B64\"}" \
"$VAULT_ADDR/v1/transit/encrypt/my-kyber22" \
| jq -r .data.ciphertext)

izi0511@Igors-MacBook-Pro vault % echo "$CT"
vault:v1:Wx1KTL9j3U80SWpxYvLbyi2uQoMQ6JHqAaBvnRYknO0eD4h9DY3iNcL39SVMVHi5QVoe8b/+lJ6alts9oAIxAFJ8QbIn5zEHJrCiw1Bu56X0BDNBErVjOFd3gE2d1KMz8+JVWwKidmDjJJxKZ6vIIL3GDcKzibHuRSrNBUsWcORnOJuls1op4AM3eRb0zgH6Zu3j6T2agtidRlo14cDeGf8JBH5Fy9KZmxVAlNPj7+BKBRXF66QczwwXy+p2selkSuqCFNoTh2zJsqL3MB06/oxT6CAeFAtKEmZJUly9y5T1lGKQjK93xNZZIqpDoXtHegUT5WTQxuTGn0ag65i9XyWUZ4f93M2VTLtab2qjOBMaX59ll5RUVG/5tmg7SNoa1eD7mF8efS668fPT+gHcxwcbKV8Wt4r6PbGeQuMoCE4cq21cmQmUiGWbaHcTt/Dlho42q4IjqIkAZxIKivhp9k6s+lftNDYqi3WlYVCO6JAPHU8I4G6i3FFngPV2iWoCu28q3v01xWXFO2QROpSnGCRLymmmjUiyU8uDzXTjy7Bsn3UwpgiuaIWMWJg+zN0Fl7eAgTFlYSe2d8lzAqXwMVfAl/ktJ9XCD+NT+oPDnMHJACLghst8siQfNkeeBh41x6uEJnXAkWTRdruA9XAXn4ZRi9SJdcwJFY4N6+RPXDTL7Uu276VQNibpN920tMT+TqLf8ERx7l2P9v1wzfaGoZ5pIw1MPTOe2H6o1/20cBetLYQiwSgCZjhTYE2DNhHBxrVSidvEU2J/6PXDDsBpQyZM8uiHLvadO+l8PBdFOexvrzhRyR8wuZHJJYU82gU8pga/NxNmNgUI1m8hSFfXYPy80iceryIhO8mTUu0Q0JgsFA4rXVuBUJxTg2LA/NkO93ZG/yQUKFOVZT5x0SiL5+ZmZmzHOoP7Sd2zJv8vhjzrIxKNHrkpCQJoIGVZpLLYGHwmovioYpARHiniGhSufNmv/r28vyZ1HOhZInvStLX35swBJ/OmkJxPT76EWRQX8a+ei0W976r4zaitfXLuTm0nb/RoMWbcNBbZw1Q+YZbOH+KM25uAC+aTm5Dgh4obFE7gFun0bhn0+1az0+d3XK+kjszIb6gNO9sNcnmVgdeDTWGnzvwFyD1HpBuSfx647S1c07eppv4BI5eL6mGkhJlvAY1k1rGeZ2OdLR8DlF0tMnJOEQI3Hddji+f/JovSMndrkQ4BB5QtsLoF4/ve3JOOI+VGWeLIlyS80e4JFzRsWxwAtxLe+NoE5I2WOWhKNwO4KNzCl21YIFax3TwVC6L3NRnwuHj5Nj8fTd9bvZkMG5HU06rHiFVfVE2G3qc8mXImUjBETzjpsp+65rn90KW7LBZdgm5FH/8M2bvOy9ENgZxX3+a4o1BD8XW3Fa7dcXYzcv1VNAtJpriBnXyFEZWyzYAUX4PIf+2GFusR25W/qLfnmN04a2Tgkp5KwEarOCCYzWac2Aj/7C+5Bdeq7Esue7gCKzmkOgg=

izi0511@Igors-MacBook-Pro vault % PT_B64=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
-X POST -d "{\"ciphertext\":\"$CT\"}" \
"$VAULT_ADDR/v1/transit/decrypt/my-kyber22" \
| jq -r .data.plaintext)

izi0511@Igors-MacBook-Pro vault % printf "%s" "$PT_B64" | base64 -D; printf "\n"
Hello Kyber768

PATH ROTATE:

izi0511@Igors-MacBook-Pro vault % PLAINTEXT_B64=$(printf 'hello, kyber 🚀' | base64)

izi0511@Igors-MacBook-Pro vault % CT=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d "{\"plaintext\":\"$PLAINTEXT_B64\"}" \
"$VAULT_ADDR/v1/transit/encrypt/my-kyber22" | jq -r .data.ciphertext)

izi0511@Igors-MacBook-Pro vault % echo "$CT"
vault:v2:36CrY7RRqcsdi8mbrKQdzbrqA/KI5XHMblIDNSrVjcE1xUOazM1+pPqYnNY52f0B3bGap49/KlryS63qEuJXEY3N8FgnAUt2AscozQV86PQMk4tFaPcizWQmd/DQkT1bD7TCvJ/0NYJGSkBLq44l4Pss0OeCn7G/PDkYrAbY8UXPYqnd5jynzr4ggaJCVSKnAmBVCk9hof5K9qwQtMnKqPn6KRw/XZ1U6DElpqoVwYtJxYWj0XV+mgMljnFIx8v9AxkWVoCQ/oi2aob5ks3ChDtj6Bmic0R4BikDGG32lRLVxtQU0awyiR91k0eHGpJzhULWjPC3/rY9eYJoLc3PADZAfYuVj9mS9dY0zfbFRtYWQxlk+ky+AJWk59DtYtVWYbnvLSvL8Z/KN9/CRmGbUKBZyXKpZSK15Pdcu8Wi8CwOgpVtZXfj6F64L/iQtLJirWY9Rd4I/8mI4AKvYOXzWzmBKU6lED3oizwfvevwLMxxXh9fIIhfQmWz+NlHVLOyCjdeMeLpkDpbwb9wrFUyZcwNeewCUQ/16e+FM8gQq3771JNYuYfcbA9Q55oKo4SRv6XWOdLa7yHNxKI8GyYNo6a3VznmvD1CV8Jl/fGBuAh4Zqnsg2fam/U2jGyyKWAlWba0MbjJmL4yn4bXUMj/tjRaueXyxjdx4xTcdXUQHf1dq+iyT9TZivxLcYSyPOhdNYzq0URt/H1wvePFvwDH8gWiY6XA2zKEdOgLmEDpY48b95/4KEHbxO+/XI+b7EsOwIbsr3Xra6hIebvIQUIq1GNQWdHmvhFLRVSXzbeges9VlCSjWw/29yg84tkRte4Qx1f9vSI9QMKxsFxuVxksepDpFERjyVlLxr4vA8snc/tLg1w911kM0lzdtgUfTQ0OS29VOgMoXd5bslgFuUwY1PHIKBKIrG143nimHCBdat8sJ+yxw5LVNJSvNlr9FGKol6zba8O9aM/1nn1bUYoDxQhl/Ux51/WOzr08OGp7SArubxCu2SZa7zi+tUX1Eswa4V55Y660NGT/Gx9FpMqIpsXTC48MORVZbvhihBmWO/fAPMeQ7Gm5ID5vmxn7/T4vlF+cOclYN/eg2+JL3NSF+CHsxXRxBJlGLVxRpbtKB2Ao1GlAdyslWI+S/SM2DCXEEh8671q+AnBApemOJvdDmvt6sISzCk/xDZdD01gNFeuj0pg5pQzQHvZPg3r4K02+nBET6vClyWRmqIMunEG6t8bcU5rI4MMyQcee6LgqBhgeMEbA5QwY+cKB/HnA0pXaVfUwJ/jliGp/KBu4IMBrbVT6/6QwAwqSbJoKjAcbMkXOgJR9D4HQ5I5FDuzpTd1MURI611Wm8nW5RIbXeMuaP+m3bxTLbXbZ1EL7QNIPMx0OK/xHcueUtX8ho0iSofKhbyU5xOYB+DOUIRYINLLkOKxsARw5pKvFUsYR+ZgTZiIo+NDKtV3pjB+lOHEfdbtJET/RuuUPtSwcuKx06uWx/kTpNdB3BaHboGX8Juw=

izi0511@Igors-MacBook-Pro vault % PT_B64=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d "{\"ciphertext\":\"$CT\"}" \
"$VAULT_ADDR/v1/transit/decrypt/my-kyber22" | jq -r .data.plaintext)

izi0511@Igors-MacBook-Pro vault % printf "%s" "$PT_B64" | base64 -D; printf "\n"
hello, kyber 🚀

izi0511@Igors-MacBook-Pro vault % curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d '{}' \
"$VAULT_ADDR/v1/transit/keys/my-kyber22/rotate"
{"request_id":"42a5f450-3669-ae0a-5c28-921e8e867366","lease_id":"","renewable":false,"lease_duration":0,"data":{"allow_plaintext_backup":false,"auto_rotate_period":0,"deletion_allowed":false,"derived":false,"exportable":false,"imported_key":false,"keys":{"1":{"certificate_chain":"","creation_time":"2025-08-21T09:18:32.732955+03:00","hybrid_public_key":"","name":"crystals-kyber","public_key":"omFqJllRBOIHKDsq+AQFlxKgiQum73cnivSIKtKE5RofORZrk3deqWWnuNJ2UsSPrzdCa+LDDANsaQcM1piRqmGHpvSqZFEgw2jCo2LK0xuSJ+XMH3O14jhpJCK9nvSEjJuJztPH1bQgveIqJiFji+qyBdZRo5UUtShOgWs/trkv/KtknGtGdwPIPKcySbwlxeRwvvySWpBjZgAQpBa9+VicGqeY9He2zVMndqZszWMWZksgHnej5MuyUfIpqvchr5AMueA4ZZJFFqc+ErVXT2TFqVFOAqqSoQDQB/OF7SmbGvhl80oNA3GOHigKhwVBKPxjRxaIdPt9SckrI7K1Fis13ty0aACLFzIsRJgBvAlK/lRh7Atra6QHm3en3ZJCQKbAtBmbboSfTgkKDhqacxwYs3fKCXJ0kVpxU/yQr+Mn5yF3tCes5HfNa1VGEunDXGSyTtYyIMylWphhTOAlHVE0Z5KD5BQ0kxlc9DGfQvnDBEuOazKFQhKqyhgq0LsEZjwLHAS7giYacJEGHxCbUVd9HFbHAvBy5yl1q/wv1ZpVsnE6DCSqkeIRwsicT/dUyeQjF+LJIum7kBMSrQRakiZQ9zyI4vKW+bSc76CelLE1wYWYySena0YfLPAwJDVvygY66JqJ8dK8fPdFKYuz+oN1cKANk4CYFUBmfyEG50syMRpwNLJ5JsNGAWu+mtggprkyssBTM0Y+0AJdw4eu01gCnoVNABy4rSN0+JIpm1UhVyofKpprShmQNaGbk4VYw/wIFUAyoaM9gnVDb7lANahzO5lz5MmKAwCJWNlvqlzC1pYRfUiPUlmh8Jx6kjSKvDISvhJk45grdNNjwgcc7eS4ZVtHSxC50LJa/UdKGjqXHsEB/MATRIQQ1YAaJXeW9/MZzEgbDqaMiPks+jYMO0JrbGcv2QFX3lWSjCaj7lHDQmLPPIGsNUNB28ukJrzNOPCO7bSvr7x80ELNLKevx0hXekRxI1F8jqeMUbt+SGzDF3q5P/ajm6EGuDzIwWUqjRkMzYZML7eiZilZFgXHYqxLBZOUwkUq59qRMzlnsZGqWIsxEikRwNpmDoFuI9ZYswFD2gA2GeUjm9eKsEFxmmWpdWymX1UMwMwJpTQN8iHBLmeh+oZc4ie9m0uqS7mfQMEMSdDGDdG56lRMhXuHMSahnPtafTImPQK1/POOI9oizcScx/SF+6muKvslZJSh34xvZqU+jyqnOcAlYHYCwPO4JEhVAnujXJZvZoyY0fZOgiOMuuFTkQFnaNorYQEzixOm4qpP4qunr+p9Zgon78K1BOk/w7aMpBpxLfa/aKlh/9hP5YR3qXUfh4mEQtuYBABZcgmFU5ZgAmOk3NNvOzwEpMLKQzpiYhpHKzAK9ssjy/fKGjkgj5yX73jEz4WjHkMJYjaFgXOvb/Im2dWPK5kDndjIrBIrY+eHnbeZvLU4tWuszvV6sqg+hKgksyJIbSAzq8xuW3a0tfSbdAK6+VUOdhAajwlOyvcDgCk+xslx6FoLEiwNTwG4Pda98OwnHOGsoLlVXDDGz/tcIUxdBhXfNg5rXxPLFCJ4cEMYmdU="},"2":{"certificate_chain":"","creation_time":"2025-08-21T09:27:24.473825+03:00","hybrid_public_key":"","name":"crystals-kyber","public_key":"H/B6wwsze4bHDqu/HrXJ0EOm0GorI3w1+ep05UIh1MhtGzZ7ikp9YiJOnlBD7qpNdRtIMmUJ3gtgCsgT45R8v4G1RVUu3oG9kbyLncpd9EnPf/pdM0U/C5wzydwHqXuEIJjNRWywdhs5ePge0TKlDoV0ImNx5QLN7AlHD9ctADYShIQytJm2htGI9Vm/JSFJoCcYXtdbSfC4iwp/S1eKx1mYI5ErEXkTcUSFaDq//ixdKKGyA1eEI2F5J8IKCyax9KSWsPY60MfAPSQcPjwqJ9KppWmfhecpiqmPNXGfOOHD1Gm//+x8AdTEjfkgITkSt4YQ1neg58ioQOU8opsWMVvB1Lq6N8OsG5KW1oamqHydJXwbOWM4IkV8QqkFv/WF7MAel1pZ1LSCNVa7p0F1nxlZeJrPCaSu7toPXYc49eDDUTaM1xJaQMpOP0MwwuoxByzN4dq7XcA6B9qLzWemY+AdTPwFRIgIpRO0eoQbZJmbSdGKUVyFEAIGg8VEX8K8EOSgzZpnhxlgm8qdcCC/ybpm68JgCPcr3SG8dFfPTtWTkYeqMViJSqpmCtpDalzPXOZYDLNnaiiLSecKsWJNuxVNNuQ8dBFyRXhszbEhChKwj8d+u6ldT1Aaymh5dnsaAed3efvHEnHMoLzKJXmK65kl2nC6E2wSF3Vc+0WwegGHqwswlVd+UdZ8raYflLE02nIaD0QfmOORFrCXvAsAr6PPwVfENddOEUIx5LYjWlIcY0lk1aBaCqsp78djjByyiAMVyrQ5fQl6tQydZgq9BzuPeNGOIQw1nuF1OIi+q2mIXDE4ibKGRrl9MmMNEAJDuCxEJsZ+JjO/EPBQPGsyTqWp/hq0TXe2Nzch2WqckBZXy1PLCVSmSIECUTHEe2vHvykmY8y3fbhtjOiSUEiQ1XvMHomrDRMtOSESewQCc4IWjSEjihFbjrytjSR8CmnDi8ElZCcault04ktyuSIvYNyeY3DAfoRjKzrKOoNsIiDONKmbJmMiU9DC3KG8HFayx7aypKlRPnavXHWimvwCuqAZ98Y7ofEvnkWZAHpDAEswgjQTBYBvnjovEcs+rJQqOEGCSdRNL5FiBgEYQ6gNqaRyawJIk7ks4WCILZnPWjhdRMBEXVK+gokn27MtPDiXYoJ2UgLKDJKXNbQaBxTHyoBhIKYqcjgP4lwiRfk6AOoy2spQJywDKtC3+llIg9q6LCR1Saa+tuSgcZqaWRStHFV2Z7S8f+o78eLJekFRi6JY8QyueyVRqaMBNogCmfELm3eFhdmrMMdpWoA2H/lAhCE1PLhpdMiqBnuyfjxwfmsEC0awaJtQHaJtIpVyvEfLyph/aIAi+slI9mm9+eOa/NtL2MoQ6XsD90OOsbZ0PqTIeyMVb8bPEbVMpfMB2ZGjJIeVvCBBIkQYMzlUC9dDD7oDKXkHFnlcXQO3QNUTGoeqTHFykRouplYY0Hwm7YrLDFsT1+fBoyR74VGyeUdg2Tqp8AkjMPme1kCAxvUWHLxo/UBx+fgmdHYuoRV0EWvNa9ZYMdx/T9EfSJnQ8XSzl4SMJUE6DroHKlr98IhaP0k="},"3":{"certificate_chain":"","creation_time":"2025-08-21T09:32:06.065947+03:00","hybrid_public_key":"","name":"crystals-kyber","public_key":"v0O2o3k/28cNcKtzOBuJvBiS1Pq0zpkjAwhTMagQl5NB14x6ffmaB0A+qskAqJQ0p+gCYyy2rNs/WVFnPYxMHUUCfAEgSiVlLewwG3wzBjzClLwT7LAhlwICFnydmfYgKxjO4mkzSzBxaSoFceSf8bKRwtN2bQZez2N+S1wigdGwu7ckREeccXeoQhmkCYpbdAQ71gN+12TMK/V/gkvGJ/NvkwEcNzO3GAg68NUy81FH8kfK84TAIRxTEAiHD/It8fstS6yzhiBYS5cgLzpBmlqvX6dx3IinlCmxwSgr3+V+v1cdUUhUZfLJKppTlkF1QyqYzNAyhJcLmAJ5WqhTxRiaR6WrgIyatVdwWpp/cFu8L7uKWTZ7N1Bi7GwUCrZMn/a/MKyDIaCR5CDCu1V1XfipLrC6IyIh4ehDp/cU/QaBUKqwmdAGLDZNLHCSZsWRpqHBysWdxHKAixallCV0xzUTQwlSO8d3rPeXwEYPr1SAyvYdzEFC7zCRiqirzby5TXyEPJmF6MWvs9EQQ3khPKO7/xUa52RAkCQmwFKnKIZ2NeFB8CHFdOibWHuNUWIdaOGBRmoV8gc1NKkYTsZCmoodo1jKmqMx+UAK5kqwrod1kMKAonoXkFtrZqy54kkqCGa4aFF9a1oIWVsM3voykKMBSCu+J1ad/uc60vjDg3jBMctq+gYtE6sibcZA8NhDU7YSCtEH7rZ+PsdPHci3vPF4ALMohXK2hjt1xHEJ1rgrDxsTkHkehRa6k2cMWuQMXFlOa4GEMVuI9jJKn8mtVxxpNjtoWRdB93c8KReIOilIIRg/NvePI1Kf2ARah8c1VbceyStk+zvOC6XDJcg20VgrV+o0DrssmBpDMIKDQpo1omq1ZsK2hXPI6WUtaHIoJqUAniPPOJG50bxAJHmdsYKw7liDtOYBBDCLkGMg7YAZ+HOV5Uaqh3XC8YGm6AZj1acas5u8GwwtvPM/c5BNHRuysdaCh2I0pnaGtuRp6hs3o8vGHBudT8sEnQeW+MKOVelyZ2wwTucaxdEzIrxJtrNRIZqHL8cyenUWYkCm3xlUbBWhwIEa6wGZq7SKrSO/A+sbqQOjrltPt6GcXIZRoxFFe2YvbwJyaXSdh/G/PEOzIfiSJPCI7tyugwSD98hoEid5OHSxApYJW/WVfgM8dFd7YIrAUft7YnxOOKJTSMxbxQGNs7Afx8DCduGRHQzO4SovP0m6VlKd8pVO0kw1CHkpXiK8GAwDdQqq4KdmDoNu0Gp+XZuO2VGRxuuLG0x/B6K86zd7DrY7WHatCjUbJJefbqAvLYIQbojJO4CepXF7KFCboeoXmxBlr3Yq4QWgCCqJzbGE23SL3TdMueuNDLufxPkDpZrKh/ya9rkfwAlW5oiWpUxyCbBq9JQlP9RNUOY7eMK9b7yBKWzBGJOctMGqU2x6n4TJlrQF8PW6LuBUWxdKm7Uw1jp0geFjfxwhLBGVeZp+VLNHiMJyU7eO2CQb67SiV6as9XbEincJ/vilrRKseTGEvUW6M+Jt8xtZlvxO9JEPwG6N+a5Ytf9YJ8IzfNzEPhU/n8XW+Na5EV4="}},"latest_version":3,"min_available_version":0,"min_decryption_version":1,"min_encryption_version":0,"name":"my-kyber22","supports_decryption":true,"supports_derivation":false,"supports_encryption":true,"supports_signing":false,"type":"kyber768"},"wrap_info":null,"warnings":null,"auth":null,"mount_type":"transit"}
izi0511@Igors-MacBook-Pro vault % curl -sS -H "X-Vault-Token: $VAULT_TOKEN" \
"$VAULT_ADDR/v1/transit/keys/my-kyber22" | jq .data.latest_version
3

izi0511@Igors-MacBook-Pro vault % CT2=$(curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d "{\"plaintext\":\"$PLAINTEXT_B64\"}" \
"$VAULT_ADDR/v1/transit/encrypt/my-kyber22" | jq -r .data.ciphertext)

izi0511@Igors-MacBook-Pro vault % echo "$CT2"
vault:v3:vrD69phUNbiG1+s3NDvBnABN7GOZcLGsJQ8MGwSpWbkC0P79nEymkqa9bguVnGUsOt0UKDFUJW0CpwO4UiJ4aROK4erCgY3MqEFekQT19Ss50UJul+vARGCoVEa7gEe3NbUxTAhp+Rdk2zQjAxuWUmx7EoXD4pTeyoyd29qbCV1Z9dOHvAQGxUp3J6n/Fx5AGgv2uesOOIJuxz02KU6nV3wP7fj3OFwudfy09gLNs3CPMVBTNXuuAJ/lt0c3+TNe/BUGr/l03DihB7/aA3Bi2FTkyS/djWfP2n+IrhiVGFgnI1bRGK1oC1ib+DQq2xhv73zDHu4S1AL5UinzFyDW640FV6HRM4tpTqx2pkNYXTHNq4TlTBVCj5IsQPwglPGWgux52CZgeJrsjFJFWQH0plzizxiPOyRgYwoiuj0WOBKKZYzPxj9JEgySyeNOjscuKOZVaso4UpUV1cvnEzyh9RnPrLQrVVvx+1GovZOPXHHdjDM5VddJ+3gQeQomJSmFOCFUdC8JDuHVDhQpIYJfwXsHIeHD8xqc0Xc05mOD+4+atEBvyMgus7TUR7WBYb6hiCuUjP1ajNr0XwmFBsl+CQvjIy8vKA/smujzDf3bXje3WeYj0emf2Ww6S8ozfO5XmDyXBm2DyoyEG63VTVhNUQ7ZVAhSfRQmTGfgK/1FDq7YQD6Jx4EcogEhPS92C/4wA7WiPT+pazVKivdjv/D4E7/ckrt3ay141LhaEHKtxy6LLVilLnW7gH+5JAIYRJEM3tHtaaTuIWMBMbisVfxYO5v2ufXd/+TlHkcNFy/gQnfVl6zuepZOVI2PVuwxHt1NhkZxcaXhdxDHg5goDeBL7rozF/uPbNSaneJc5bYpeRDeWlRA9I40jRaMSw8wv57J9iEe3/1rgWTIeh1OkZTAlIgKLUV+Bl/bblii6TB3js3hg+qYTC+b7NkwEt8ht3PXxy70pgsibRFkWkJq0OPSZvzo7eG3wrDp90GJntmcYn6wh0PODOKGJD3dmzCQkpxhE3c9hycQVUPVcHkKaT+ygBGuCVMPGlWF3xATaHTlqY4eFBcd/JOeLDDfErP5Z+CW/JL8Tsg95z/wyykaSZdADmWGZ5rRj92jS8jgAQRTOSeNVZVsXvCzbaqoZuFJnYSS+OrrK0XxO0OhhmWQ/nrUNW/lNTEfr2w9FTJuCJGQW2ZAve7Jx93YN9s8sTDeqB683P6MytXRcrtr+K3m6ako0Tgkyd3B1M0XW+4WgEReYcY6WRDWZHNjI1/qkv7Jxh/+NbBb+4a3FEHyhq9blokpaZV1IIbtK28AKSCYVzWql3YUXtQ2NUtuqycVq7s/SY05JNCNUdUgWr9rvbr6XUJT4lh7xCyoUGZpK3gAvRTMn8e8jO1xloyN78rtAJw0syDkyeHaOTLGMB8Y3ZrxJ119wzSwxNHpOg40iQ+m1mo5YYwkD9QapmDFtXO7SQr1oDgaUnO7FQ8SbWfdaLkVuilDB780Rk8OjXLaYf88pp4=

izi0511@Igors-MacBook-Pro vault % curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d "{\"ciphertext\":\"$CT\"}" \
"$VAULT_ADDR/v1/transit/decrypt/my-kyber22" | jq -r .data.plaintext | base64 -D; printf "\n"
hello, kyber 🚀

NON HAPPY PATHS:
curl -sS -H "X-Vault-Token: $VAULT_TOKEN" -H "Content-Type: application/json" \
-X POST -d "{\"plaintext\":\"$(printf 'hi'|base64)\",\"nonce\":\"$(openssl rand -base64 12)\"}" \
"$VAULT_ADDR/v1/transit/encrypt/my-kyber22"

{"errors":["provided nonce not allowed for this key"]}

