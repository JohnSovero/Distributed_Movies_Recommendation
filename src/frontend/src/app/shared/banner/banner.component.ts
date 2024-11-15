import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { environment } from '../../../environments/environment';
import { WebsocketService } from '../../core/services/websocket.service';
import { HttpClient } from '@angular/common/http';
import { Movie } from '../../core/models/movie.model';
import { CommonModule } from '@angular/common';
import { LoaderComponent } from '../loader/loader.component';

@Component({
  selector: 'app-banner',
  standalone: true,
  imports: [CommonModule, LoaderComponent],
  templateUrl: './banner.component.html',
  styleUrls: ['./banner.component.css']
})
export class BannerComponent {
  bannerRecommendations: Movie[] = [];
  private bannerSubscription!: Subscription;
  currentMovieIndex: number = 0;
  showRecommendations: boolean = false;

  constructor(
    private websocketService: WebsocketService,
    private httpClient: HttpClient
  ) {}

  ngOnInit(): void {
    this.connect();
  }

  connect() {
    this.websocketService.connect('ws://localhost:9015/recommendations/above-average');

    this.bannerSubscription = this.websocketService.messages$.subscribe(
      (data) => {
        this.showRecommendations = false;
        this.bannerRecommendations = data;
        console.log(this.bannerRecommendations);
        this.startBannerRotation();
        this.showRecommendations = true;
      }
    );
  }

  startBannerRotation() {
    setInterval(() => {
      this.currentMovieIndex = (this.currentMovieIndex + 1) % this.bannerRecommendations.length;
    }, 15000);
  }

  ngOnDestroy() {
    this.bannerSubscription.unsubscribe();
    this.websocketService.disconnect();
  }

  prevSlide() {
    this.currentMovieIndex = 
      (this.currentMovieIndex > 0) ? this.currentMovieIndex - 1 : this.bannerRecommendations.length - 1;
  }

  nextSlide() {
    this.currentMovieIndex = 
      (this.currentMovieIndex < this.bannerRecommendations.length - 1) ? this.currentMovieIndex + 1 : 0;
  }
}