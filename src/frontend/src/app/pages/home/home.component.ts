import { WebsocketService } from './../../core/services/websocket.service';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { Movie } from '../../core/models/movie.model';
import { Subscription } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent implements OnInit, OnDestroy {
  bannerRecommendations: Movie[] = [];
  private bannerSubscription !: Subscription;

  constructor(
    private websocketService: WebsocketService
  ) {}

  ngOnInit(): void {
    this.connect();
  }

  connect() {
    this.websocketService.connect('ws://localhost:9015/recommendations/above-average');

    this.bannerSubscription = this.websocketService.messages$.subscribe(
      (data) => this.bannerRecommendations = data
    );
  }

  ngOnDestroy() {
    this.bannerSubscription.unsubscribe();
    this.websocketService.disconnect();
  }

}
