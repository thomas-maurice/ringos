#include <auth.h>

bool needsAuth(String password) {
    return password.length() != 0;
}