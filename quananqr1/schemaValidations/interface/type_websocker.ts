export interface WebSocketMessage {
    type: string;
    content: any;
    sender: string;
    timestamp: string;
    tableID?: string;
    orderID?: string;
  }
  
  export interface OrderContent {
    orderID: string;
    tableNumber: string;
    status: string;
    timestamp: string;
    orderData: any;
  }