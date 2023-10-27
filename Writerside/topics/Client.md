# Client

NewClient
 : Creates a new OpenSky client.
    ```go
    func NewClient(opts ...ClientOption) (*Client, error) {}
    ```

GetStates
 : Gets flight states using the background context.
    ```Go
    func (cli *Client) GetStates(bbox *BoundingBox) (StateResponse, error) {}
    ```

GetStatesWithContext
 : Gets flight states using the provided context.
    ```go
    func (cli *Client) GetStatesWithContext(ctx context.Context, bbox *BoundingBox) (StateResponse, error) {}
    ```

