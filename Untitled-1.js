const socket = new WebSocket('ws://ec2-54-153-173-231.ap-southeast-2.compute.amazonaws.com:8080/ws');

socket.onopen = function(event) {
    console.log('WebSocket connection established');
};

socket.onerror = function(error) {
    console.error('WebSocket Error: ', error);
};

socket.onmessage = function(event) {
    console.log('WebSocket message received: ', event.data);
};
