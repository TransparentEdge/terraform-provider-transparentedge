---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_custom_certificate Resource - transparentedge"
subcategory: ""
description: |-
  Manages Custom Certificates
---

# transparentedge_custom_certificate (Resource)

Manages Custom Certificates

## Example Usage

```terraform
resource "transparentedge_custom_certificate" "mysite" {
  publickey  = <<EOF
-----BEGIN CERTIFICATE-----
MIIFejCCBGKgAwIBAgISA0fNiEs0uD/H6tAPWIJoL96cMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yMzAxMjYxMTAxNTBaFw0yMzA0MjYxMTAxNDlaMCoxKDAmBgNVBAMT
H3N0b3JhZ2UuZGVtby50cmFuc3BhcmVudGVkZ2UuZXUwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDPNFwUuZgucJI8lC5bEdAlfB+fIqpbBe5IVnPzTk7r
j+lxsp6aP6uIEfukCQYBWKcIk96SNiALSYrgRPnfMtyrU0rBNUHYit6b0XOmp9OY
wLa20gqhg0b2NQiQIM3WQqFTnSK+O5kvFFqoK4CY3L2cczk8GwJrMjKnabBWWtv8
hECgpztuF8fWoP/qO1OfWsMpfhBimP90JuUmKV/Qn7kkjD+m/Xb7beK6oYL7mq+B
ikTPTKNZeaO++88nFuIy5tdWwP9/km74P026mqXFNMI6wpfFDA0Q9d5QSt5152E2
HipOLv9tYvw4O+UsvsvIgTD2aErA4YtIA6KnGqrPEOXNAgMBAAGjggKQMIICjDAO
BgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwG
A1UdEwEB/wQCMAAwHQYDVR0OBBYEFN6Rzy9Gr44X9xO9e6Zrzn8s7QV+MB8GA1Ud
IwQYMBaAFBQusxe3WFbLrlAJQOYfr52LFMLGMFUGCCsGAQUFBwEBBEkwRzAhBggr
BgEFBQcwAYYVaHR0cDovL3IzLm8ubGVuY3Iub3JnMCIGCCsGAQUFBzAChhZodHRw
Oi8vcjMuaS5sZW5jci5vcmcvMGAGA1UdEQRZMFeCG2FsdC5kZW1vLnRyYW5zcGFy
ZW50ZWRnZS5ldYIXZGVtby50cmFuc3BhcmVudGVkZ2UuZXWCH3N0b3JhZ2UuZGVt
by50cmFuc3BhcmVudGVkZ2UuZXUwTAYDVR0gBEUwQzAIBgZngQwBAgEwNwYLKwYB
BAGC3xMBAQEwKDAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5sZXRzZW5jcnlwdC5v
cmcwggEEBgorBgEEAdZ5AgQCBIH1BIHyAPAAdgB6MoxU2LcttiDqOOBSHumEFnAy
E4VNO9IrwTpXo1LrUgAAAYXt9KHoAAAEAwBHMEUCIQCZeppNC9bXVXwW/9Cn4tQ8
bgRxTf69D395hH+Zyl6w8AIgIVp+P8wE4ObRh+wm1duEbjRzb9Uyu4Kev/2i/X2Q
lNkAdgC3Pvsk35xNunXyOcW6WPRsXfxCz3qfNcSeHQmBJe20mQAAAYXt9KIEAAAE
AwBHMEUCIQDdE0KH59uZM1uoBjg+X5ZyPwr1T3Cau8SXbcIueIFCOQIgWHfT7R9P
9SNk19DoSO/vhJIa8cj7yNYapuQ+AQ2GkHwwDQYJKoZIhvcNAQELBQADggEBAB0W
vkTTIXP2hZbgNxypVWXcqVNSQCj9UV90OlS5/bIMtBRInsl5wEPEnlAKpTPQWDDs
lvoHDMvt9ZqE+UpjxI1/FKKTMb4vTBLBObzrTiYSyuWv3vdwYWa3r3VnR52iAKIO
ZP9lTyH19vVjoeD3tJB3rRCm4LFlfj3HIacjFaEAb59ujnb1+Y23B4tMDzOBQ+lM
BSEuNDwRzus7AYPII+klo6DBeLTdMq8VwR1rNqdYGNuTGcdV4YeP4Y42pEE3UWp8
iV6N+VUgxWaCchVR9/vQBtUfGRMIUIXxsjbur+LquLiW1fr550oMsxIER1vA5gxv
n+lpd31uABr3m7enMFM=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjAwOTA0MDAwMDAw
WhcNMjUwOTE1MTYwMDAwWjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDELMAkGA1UEAxMCUjMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC7AhUozPaglNMPEuyNVZLD+ILxmaZ6QoinXSaqtSu5xUyxr45r+XXIo9cP
R5QUVTVXjJ6oojkZ9YI8QqlObvU7wy7bjcCwXPNZOOftz2nwWgsbvsCUJCWH+jdx
sxPnHKzhm+/b5DtFUkWWqcFTzjTIUu61ru2P3mBw4qVUq7ZtDpelQDRrK9O8Zutm
NHz6a4uPVymZ+DAXXbpyb/uBxa3Shlg9F8fnCbvxK/eG3MHacV3URuPMrSXBiLxg
Z3Vms/EY96Jc5lP/Ooi2R6X/ExjqmAl3P51T+c8B5fWmcBcUr2Ok/5mzk53cU6cG
/kiFHaFpriV1uxPMUgP17VGhi9sVAgMBAAGjggEIMIIBBDAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYB
Af8CAQAwHQYDVR0OBBYEFBQusxe3WFbLrlAJQOYfr52LFMLGMB8GA1UdIwQYMBaA
FHm0WeZ7tuXkAXOACIjIGlj26ZtuMDIGCCsGAQUFBwEBBCYwJDAiBggrBgEFBQcw
AoYWaHR0cDovL3gxLmkubGVuY3Iub3JnLzAnBgNVHR8EIDAeMBygGqAYhhZodHRw
Oi8veDEuYy5sZW5jci5vcmcvMCIGA1UdIAQbMBkwCAYGZ4EMAQIBMA0GCysGAQQB
gt8TAQEBMA0GCSqGSIb3DQEBCwUAA4ICAQCFyk5HPqP3hUSFvNVneLKYY611TR6W
PTNlclQtgaDqw+34IL9fzLdwALduO/ZelN7kIJ+m74uyA+eitRY8kc607TkC53wl
ikfmZW4/RvTZ8M6UK+5UzhK8jCdLuMGYL6KvzXGRSgi3yLgjewQtCPkIVz6D2QQz
CkcheAmCJ8MqyJu5zlzyZMjAvnnAT45tRAxekrsu94sQ4egdRCnbWSDtY7kh+BIm
lJNXoB1lBMEKIq4QDUOXoRgffuDghje1WrG9ML+Hbisq/yFOGwXD9RiX8F6sw6W4
avAuvDszue5L3sz85K+EC4Y/wFVDNvZo4TYXao6Z0f+lQKc0t8DQYzk1OXVu8rp2
yJMC6alLbBfODALZvYH7n7do1AZls4I9d1P4jnkDrQoxB3UqQ9hVl3LEKQ73xF1O
yK5GhDDX8oVfGKF5u+decIsH4YaTw7mP3GFxJSqv3+0lUFJoi5Lc5da149p90Ids
hCExroL1+7mryIkXPeFM5TgO9r0rvZaBFOvV2z0gp35Z0+L4WPlbuEjN/lxPFin+
HlUjr8gRsI3qfJOQFy/9rKIJR0Y/8Omwt/8oTWgy1mdeHmmjk7j1nYsvC9JSQ6Zv
MldlTTKB3zhThV1+XWYp6rjd5JW1zbVWEkLNxE7GJThEUG3szgBVGP7pSWTUTsqX
nLRbwHOoq7hHwg==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFYDCCBEigAwIBAgIQQAF3ITfU6UK47naqPGQKtzANBgkqhkiG9w0BAQsFADA/
MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
DkRTVCBSb290IENBIFgzMB4XDTIxMDEyMDE5MTQwM1oXDTI0MDkzMDE4MTQwM1ow
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQCt6CRz9BQ385ueK1coHIe+3LffOJCMbjzmV6B493XC
ov71am72AE8o295ohmxEk7axY/0UEmu/H9LqMZshftEzPLpI9d1537O4/xLxIZpL
wYqGcWlKZmZsj348cL+tKSIG8+TA5oCu4kuPt5l+lAOf00eXfJlII1PoOK5PCm+D
LtFJV4yAdLbaL9A4jXsDcCEbdfIwPPqPrt3aY6vrFk/CjhFLfs8L6P+1dy70sntK
4EwSJQxwjQMpoOFTJOwT2e4ZvxCzSow/iaNhUd6shweU9GNx7C7ib1uYgeGJXDR5
bHbvO5BieebbpJovJsXQEOEO3tkQjhb7t/eo98flAgeYjzYIlefiN5YNNnWe+w5y
sR2bvAP5SQXYgd0FtCrWQemsAXaVCg/Y39W9Eh81LygXbNKYwagJZHduRze6zqxZ
Xmidf3LWicUGQSk+WT7dJvUkyRGnWqNMQB9GoZm1pzpRboY7nn1ypxIFeFntPlF4
FQsDj43QLwWyPntKHEtzBRL8xurgUBN8Q5N0s8p0544fAQjQMNRbcTa0B7rBMDBc
SLeCO5imfWCKoqMpgsy6vYMEG6KDA0Gh1gXxG8K28Kh8hjtGqEgqiNx2mna/H2ql
PRmP6zjzZN7IKw0KKP/32+IVQtQi0Cdd4Xn+GOdwiK1O5tmLOsbdJ1Fu/7xk9TND
TwIDAQABo4IBRjCCAUIwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYw
SwYIKwYBBQUHAQEEPzA9MDsGCCsGAQUFBzAChi9odHRwOi8vYXBwcy5pZGVudHJ1
c3QuY29tL3Jvb3RzL2RzdHJvb3RjYXgzLnA3YzAfBgNVHSMEGDAWgBTEp7Gkeyxx
+tvhS5B1/8QVYIWJEDBUBgNVHSAETTBLMAgGBmeBDAECATA/BgsrBgEEAYLfEwEB
ATAwMC4GCCsGAQUFBwIBFiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQu
b3JnMDwGA1UdHwQ1MDMwMaAvoC2GK2h0dHA6Ly9jcmwuaWRlbnRydXN0LmNvbS9E
U1RST09UQ0FYM0NSTC5jcmwwHQYDVR0OBBYEFHm0WeZ7tuXkAXOACIjIGlj26Ztu
MA0GCSqGSIb3DQEBCwUAA4IBAQAKcwBslm7/DlLQrt2M51oGrS+o44+/yQoDFVDC
5WxCu2+b9LRPwkSICHXM6webFGJueN7sJ7o5XPWioW5WlHAQU7G75K/QosMrAdSW
9MUgNTP52GE24HGNtLi1qoJFlcDyqSMo59ahy2cI2qBDLKobkx/J3vWraV0T9VuG
WCLKTVXkcGdtwlfFRjlBz4pYg1htmf5X6DYO8A4jqv2Il9DjXA6USbW1FzXSLr9O
he8Y4IWS6wY7bCkjCWDcRQJMEhg76fsO3txE+FiYruq9RUWhiF1myv4Q6W+CyBFC
Dfvp7OOGAN6dEOM4+qR9sdjoSYKEBpsr6GtPAQw4dy753ec5
-----END CERTIFICATE-----
EOF
  privatekey = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDPNFwUuZgucJI8
lC5bEdAlfB+fIqpbBe5IVnPzTk7rj+lxsp6aP6uIEfukCQYBWKcIk96SNiALSYrg
RPnfMtyrU0rBNUHYit6b0XOmp9OYwLa20gqhg0b2NQiQIM3WQqFTnSK+O5kvFFqo
K4CY3L2cczk8GwJrMjKnabBWWtv8hECgpztuF8fWoP/qO1OfWsMpfhBimP90JuUm
KV/Qn7kkjD+m/Xb7beK6oYL7mq+BikTPTKNZeaO++88nFuIy5tdWwP9/km74P026
mqXFNMI6wpfFDA0Q9d5QSt5152E2HipOLv9tYvw4O+UsvsvIgTD2aErA4YtIA6Kn
GqrPEOXNAgMBAAECggEAPcQN7t+kTbOg5A4IA3273nCxvG5I+fk6nrWmutCNFgtA
O3RTcwenylgR+0P1VlFm+Vea8VrREoxJqbDmC3LN9QRPNGj7x+EdmrVFFFjS6qYH
0Verc5n+fUYx10TwFv6luJcO1EZP04jtvVO6cdbbbteqKBClF+9OyjjnJ9bN3OfF
/ODsWhy1VK+0GtAAb2WKz7D43GEgAdVjywMJomVbrpQRx1v8V74c8VmBBMZgYi+V
tntwPSsSza//UVXz2S5TMjN+1GQgLstpSWJr4+Xer9272k6EOHvTQI3uhpNNXLGk
nz4jH3icJEGOiqVo3eg3P9+IapAE7B5+h9KfxmNSAQKBgQD2WIne38Jh8iQWlWAW
n9pybG0Fny+XdZpv/q9ZG9HQH00+Q7mtPOWUcUNv+fzlklyvATJdQ+yiZXIx8pmL
81CwV29I5GQihVSBnKzgyb5vmDvqvdcxyCMoMx6o/D6AB2J6OVfnseq70IrNwGJ7
ba4o5K3UgvCSZWeJ8/TSCDFxhQKBgQDXUyLmyi3HXjyjlcxqFXAJmALh6dvuJAIB
ag90tXJP6f426Xqr7E/YQK/P77KwCIw/+RGW5ttriOKNCx3fK6PT9qKhU6WJofCP
4WReywG3MtJrt5ucgM2GkU+neTdmUcKW6Pv2bVImcp/cxF4nou9ExBBQyhibAZhk
g/VQljCxqQKBgQDMfKhNYk5XwYklWe9+OEk7jDdfYElAH3YIG1Bw1n/uk90pn0xE
unUUKITDMa803a6j8oldE+Ic17rYLTo6Cspi5uFQj41zflusj2KN4cl7ltG9xMIZ
57kPSIfd3C0BV5/uNyV6BZ0FNFHUAyt8q4nTFigZbGvICfbNc704j2aDhQKBgFoI
piBQS4IAcmSIP1fgLN+mExZ5XX+eyMPkoB/RusGVerllOOjoP56RtbHBbTrT6Cjb
sTIix36YVvpYup3VNoRrrSa9vgrljpvqx7gnNElw07E8rbFr3gQ1gFPriHGdIDtP
ogMxRNdUuGlsJl52b4uWW6gcSNuPeDQXRRz0H9o5AoGBAK28urLKT3NsAaCEWlpQ
OdLCB6Y2H0QiKOCBgZTes1W1vbjktrhQtpF/4kJxTnCQkiCUG7z3Bjk8Rp3zgpwI
uQwBTY/IR6zz9U8VdAd52qFZ0lbkb5lW3aQU14/faDZvns2zgVBqXSJn/C8Io80I
LNdUnJzoYY+RmfURRRyvHije
-----END PRIVATE KEY-----

EOF
}

output "mysite" {
  value = transparentedge_custom_certificate.mysite
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `privatekey` (String) Private part of the certificate in PEM format, the certificate can't be protected with a password
- `publickey` (String) Public part of the certificate in PEM format, it's recommended to include the full chain

### Read-Only

- `commonname` (String) CN (_Common Name_) of the certificate
- `domains` (String) SAN (_Subject Alternative Name_) domains included in the certificate, including the Common Name
- `expiration` (String) Date when the certificate will expire
- `id` (Number) ID of the Custom Certificate

## Import

Import is supported using the following syntax:

```shell
# Import a custom certificate by its ID
terraform import 'transparentedge_custom_certificate.mysite' 321
```