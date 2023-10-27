# Client Options

Client options are used to configure the OpenSky client.

ClientOption
 : Interface definition all client options implement.

    ```go
    type ClientOption interface {
        apply(*Client) error
    }
    ```

Authorization
 : Option used to set a username/password in calls to the OpenSky Network API to access
    own state vectors and increased rate limits.

    ```go
    type Authentication struct {
        ClientOption
        username string
        password string
    }
    ```

    To create an Authorization client option, use:

    ```go
    func WithAuthentication(username, password string) *Authentication {}
    ```

Timeout
 : Used to set the network timeout in calls to the OpenSky Network API. This allows calls to return larger
    datasets, for instance all active flights instead of flights in a small area.

    ```go
    type Timeout struct {
        ClientOption
        timeout time.Duration
    }
    ```
     
    To create a Timeout client option, use:
    ```go
    func WithTimeout(timeout time.Duration) *Timeout {}
    ```

    
