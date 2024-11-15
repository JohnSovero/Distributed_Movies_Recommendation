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
  selectedIndex: number = 0;
  showRecommendations: boolean = false;
  private rotationTimeout!: ReturnType<typeof setTimeout>;

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
    const rotateBanner = () => {
      if (this.bannerRecommendations.length > 0) {
        this.selectedIndex = (this.bannerRecommendations.length + this.selectedIndex + 1) % this.bannerRecommendations.length;
        console.log(this.selectedIndex);
      }

      // Schedule the next rotation with a dynamically adjusted interval
      const nextInterval = 15000;
      this.rotationTimeout = setTimeout(rotateBanner, nextInterval);
    };

    rotateBanner();
  }

  ngOnDestroy() {
    clearTimeout(this.rotationTimeout);
    this.bannerSubscription.unsubscribe();
    this.websocketService.disconnect();
  }
}