import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router'
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { ApiError } from '../api/client'

interface LoginFormValues {
  email: string
  password: string
}

export function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [serverError, setServerError] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>()

  async function onSubmit(values: LoginFormValues) {
    setServerError(null)
    try {
      await login(values.email, values.password)
      navigate('/')
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        setServerError('メールアドレスまたはパスワードが正しくありません')
        return
      }
      setServerError('エラーが発生しました。時間をおいて再度お試しください')
    }
  }

  return (
    <div className="bg-surface rounded-xl p-6 w-80">
      <h1 className="text-xl font-bold text-ink mb-4">ログイン</h1>
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
          ログイン
        </button>
      </form>
      <p className="text-sm text-ink-muted mt-4">
        アカウントをお持ちでない方は{' '}
        <Link to="/signup" className="text-primary">
          サインアップ
        </Link>
      </p>
    </div>
  )
}
