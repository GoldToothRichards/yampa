PASSWORD=$(base64 < /dev/urandom | head -c15); echo "$PASSWORD"; echo -n "$PASSWORD" | sha256sum | tr -d '-'