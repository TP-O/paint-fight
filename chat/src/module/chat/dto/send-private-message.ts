import { IsNotEmpty, IsString } from 'class-validator';

export class SendPrivateMessageDto {
  @IsString()
  @IsNotEmpty()
  receiverId!: string;

  @IsString()
  @IsNotEmpty()
  content!: string;
}
