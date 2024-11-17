import { Component, inject, Inject, Input } from '@angular/core';
import { Movie } from '../../core/models/movie.model';
import { CommonModule } from '@angular/common';
import { MAT_BOTTOM_SHEET_DATA } from '@angular/material/bottom-sheet';
import { MatChipsModule } from '@angular/material/chips';

@Component({
  selector: 'app-bottom-sheet',
  standalone: true,
  imports: [CommonModule, MatChipsModule],
  templateUrl: './bottom-sheet.component.html',
  styleUrl: './bottom-sheet.component.css'
})
export class BottomSheetComponent {
  constructor(@Inject(MAT_BOTTOM_SHEET_DATA) public data: Movie) {}
}
