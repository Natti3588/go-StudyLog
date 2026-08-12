import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router'
import { useState } from 'react'
import { signup } from '../api/auth'
import { ApiError } from '../api/client'
import { useAuth } from '../context/AuthContext'

interface SignupFormValues {
  email: string
  password: string
}

export function Signup() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [serverError, setServerError] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<SignupFormValues>()

  async function onSubmit(values: SignupFormValues) {
    setServerError(null)
    try {
      await signup(values.email, values.password)
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setServerError('このメールアドレスは既に登録されています')
        return
      }
      setServerError('エラーが発生しました。時間をおいて再度お試しください')
      return
    }

    // POST /signup は認証Cookieを発行しないため、続けてログインしてダッシュボードへ進む
    try {
      await login(values.email, values.password)
      navigate('/', { replace: true })
    } catch {
      navigate('/login', { replace: true })
    }
  }

  return (
    <div className="bg-surface rounded-xl p-6 w-80">
      <h1 className="text-xl font-bold text-ink mb-4">サインアップ</h1>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <label className="block mb-3">
          <span className="block text-sm text-ink-secondary mb-1">メールアドレス</span>
          <input
            type="email"
            className="w-full border border-hairline rounded-xs px-2 py-1.5"
            {...register('email', { required: 'メールアドレスを入力してください' })}
          />
          {errors.email && <p className="text-error text-sm mt-1">{errors.email.message}</p>}
        </label>
        <label className="block mb-3">
          <span className="block text-sm text-ink-secondary mb-1">パスワード</span>
          <input
            type="password"
            className="w-full border border-hairline rounded-xs px-2 py-1.5"
            {...register('password', { required: 'パスワードを入力してください' })}
          />
          {errors.password && <p className="text-error text-sm mt-1">{errors.password.message}</p>}
        </label>
        {serverError && <p className="text-error text-sm mb-3">{serverError}</p>}
        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-primary text-white rounded-md py-2 font-medium disabled:opacity-50"
        >
          サインアップ
        </button>
      </form>
      <p className="text-sm text-ink-muted mt-4">
        アカウントをお持ちの方は{' '}
        <Link to="/login" className="text-primary">
          ログイン
        </Link>
      </p>
    </div>
  )
}
