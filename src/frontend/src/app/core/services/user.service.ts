import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private userId: number;

  constructor() {
    // Generate a random userId between 1 and 10
    this.userId = Math.floor(Math.random() * 10) + 1;
  }

  getUserId(): number {
    return this.userId;
  }
}