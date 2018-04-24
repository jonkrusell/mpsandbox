import { Injectable } from '@angular/core';

import { Observable, Subject } from 'rxjs/Rx';
import { WebsocketService } from './websocket.service';

const CHAT_URL = 'ws://192.168.10.171:8081/ws';

export interface Message {
	type: string,
	value: string
}

@Injectable()
export class InhabitantsService {

public messages: Subject<Message>;

    constructor(wsService: WebsocketService) {
        this.messages = <Subject<Message>>wsService
        .connect(CHAT_URL)
        .map((response: MessageEvent): Message => {
            let data = JSON.parse(response.data);
            return data;
        });
    }
}
